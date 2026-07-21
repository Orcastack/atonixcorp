from __future__ import annotations

from calendar import monthrange
from datetime import date, timedelta
from decimal import Decimal, ROUND_DOWN
from typing import Iterable

from django.contrib.auth.models import User
from django.db import transaction
from django.db.models import Sum
from django.utils.dateparse import parse_date
from django.utils import timezone

from atonixcorp.models import ChartOfAccounts, JournalEntry, TaxCalculation, UserProfile, GeneralLedger

from .automation import (
    notify_certificate_released,
    notify_exercise_status,
    notify_exercise_submitted,
    notify_upcoming_vesting,
    notify_vesting_milestone,
    sync_exercise_payment,
    sync_payroll_tax_event,
)

from .models import (
    AccelerationType,
    ApprovalDecisionStatus,
    EquityExerciseApproval,
    EquityExerciseRequest,
    EquityGrant,
    EquityGrantType,
    EquityHolding,
    EquityPayrollTaxEvent,
    EquityShareCertificate,
    EquityTransaction,
    EquityValuation,
    EquityVestingEvent,
    ExercisePaymentMethod,
    ExerciseRequestStatus,
    GrantLifecycleStatus,
    PaymentStatus,
    TerminationTreatment,
    TransactionType,
    VestingEventStatus,
    VestingEventType,
    VestingInterval,
    CertificateStatus,
)


_INTERVAL_TO_MONTHS = {
    VestingInterval.MONTHLY: 1,
    VestingInterval.QUARTERLY: 3,
    VestingInterval.ANNUAL: 12,
}


def add_months(base_date: date, months: int) -> date:
    base_date = ensure_date(base_date)
    month_index = base_date.month - 1 + months
    year = base_date.year + month_index // 12
    month = month_index % 12 + 1
    day = min(base_date.day, monthrange(year, month)[1])
    return date(year, month, day)


def quantize_money(value: Decimal | int | float) -> Decimal:
    return Decimal(value).quantize(Decimal('0.01'))


def ensure_date(value) -> date:
    if isinstance(value, date):
        return value
    parsed = parse_date(str(value))
    if parsed is None:
        raise ValueError(f'Invalid date value: {value}')
    return parsed


def _decimal_percent(value: Decimal | int | float) -> Decimal:
    return Decimal(value) / Decimal('100')


def _split_units(total_units: int, periods: int) -> list[int]:
    if periods <= 0:
        return []
    base_units = total_units // periods
    remainder = total_units % periods
    schedule = [base_units for _ in range(periods)]
    for index in range(remainder):
        schedule[-(index + 1)] += 1
    return schedule


@transaction.atomic
def rebuild_vesting_schedule(grant: EquityGrant) -> list[EquityVestingEvent]:
    grant.vesting_events.exclude(status__in=[VestingEventStatus.EXERCISED, VestingEventStatus.FORFEITED]).delete()

    if grant.vesting_interval == VestingInterval.CUSTOM and grant.custom_schedule:
        created = []
        for index, item in enumerate(grant.custom_schedule):
            created.append(
                EquityVestingEvent.objects.create(
                    workspace=grant.workspace,
                    grant=grant,
                    event_type=VestingEventType.MANUAL,
                    status=VestingEventStatus.PENDING,
                    vest_date=ensure_date(item['date']),
                    units=int(item['units']),
                    source_reference=f'custom-{index + 1}',
                    notes=item.get('notes', ''),
                )
            )
        return created

    interval_months = _INTERVAL_TO_MONTHS.get(grant.vesting_interval, 1)
    total_periods = max(grant.vesting_months // interval_months, 1)
    units_per_period = _split_units(grant.total_units, total_periods)
    cliff_release_date = add_months(grant.vesting_start_date, grant.cliff_months)
    cliff_units = 0
    created = []

    for period_index, units in enumerate(units_per_period, start=1):
        vest_months = period_index * interval_months
        vest_date = add_months(grant.vesting_start_date, vest_months)
        if grant.cliff_months and vest_months <= grant.cliff_months:
            cliff_units += units
            continue
        created.append(
            EquityVestingEvent.objects.create(
                workspace=grant.workspace,
                grant=grant,
                event_type=VestingEventType.SCHEDULED,
                status=VestingEventStatus.PENDING,
                vest_date=vest_date,
                units=units,
                source_reference=f'schedule-{period_index}',
            )
        )

    if cliff_units:
        created.insert(
            0,
            EquityVestingEvent.objects.create(
                workspace=grant.workspace,
                grant=grant,
                event_type=VestingEventType.CLIFF_RELEASE,
                status=VestingEventStatus.PENDING,
                vest_date=cliff_release_date,
                units=cliff_units,
                source_reference='cliff-release',
            ),
        )

    return created


@transaction.atomic
def apply_termination(grant: EquityGrant, termination_date: date) -> EquityGrant:
    termination_date = ensure_date(termination_date)
    grant.termination_date = termination_date
    if grant.termination_treatment == TerminationTreatment.FORFEIT_UNVESTED:
        future_events = grant.vesting_events.filter(vest_date__gt=termination_date, status=VestingEventStatus.PENDING)
        forfeited_units = sum(event.units for event in future_events)
        future_events.update(status=VestingEventStatus.FORFEITED)
        if forfeited_units:
            grant.forfeited_units += forfeited_units
            EquityVestingEvent.objects.create(
                workspace=grant.workspace,
                grant=grant,
                event_type=VestingEventType.FORFEITURE,
                status=VestingEventStatus.FORFEITED,
                vest_date=termination_date,
                units=forfeited_units,
                source_reference='termination',
                notes='Automatic forfeiture on termination.',
            )
    elif grant.termination_treatment == TerminationTreatment.FULL_ACCELERATION:
        apply_acceleration(grant, 'termination-full', termination_date, Decimal('100'))
    elif grant.termination_treatment == TerminationTreatment.ACCELERATE_TO_CLIFF:
        cliff_units = (
            grant.vesting_events.filter(vest_date__lte=add_months(grant.vesting_start_date, grant.cliff_months), status=VestingEventStatus.PENDING)
            .aggregate(total=Sum('units'))
            .get('total')
            or 0
        )
        if cliff_units:
            apply_acceleration(grant, 'termination-cliff', termination_date, Decimal('100'), limit_units=int(cliff_units))
    grant.save(update_fields=['termination_date', 'forfeited_units', 'updated_at'])
    return grant


@transaction.atomic
def apply_acceleration(
    grant: EquityGrant,
    trigger_name: str,
    trigger_date: date,
    percent: Decimal | None = None,
    limit_units: int | None = None,
) -> EquityVestingEvent | None:
    trigger_date = ensure_date(trigger_date)
    summary = calculate_grant_summary(grant, as_of=trigger_date)
    remaining_units = summary['unvested_units']
    if limit_units is not None:
        remaining_units = min(remaining_units, limit_units)
    if remaining_units <= 0:
        return None

    acceleration_percent = percent
    if acceleration_percent is None:
        if grant.acceleration_type == AccelerationType.DOUBLE:
            acceleration_percent = Decimal(grant.double_trigger_acceleration_percent)
        else:
            acceleration_percent = Decimal(grant.single_trigger_acceleration_percent)

    accelerated_units = int((Decimal(remaining_units) * _decimal_percent(acceleration_percent)).quantize(Decimal('1'), rounding=ROUND_DOWN))
    if accelerated_units <= 0:
        return None

    event = EquityVestingEvent.objects.create(
        workspace=grant.workspace,
        grant=grant,
        event_type=VestingEventType.ACCELERATION,
        status=VestingEventStatus.VESTED,
        vest_date=trigger_date,
        units=accelerated_units,
        trigger_name=trigger_name,
        source_reference=f'acceleration-{trigger_name}',
    )
    return event


@transaction.atomic
def mark_vesting_events_as_of(grant: EquityGrant, as_of: date | None = None) -> None:
    as_of = ensure_date(as_of or timezone.now().date())
    pending_events = list(grant.vesting_events.filter(status=VestingEventStatus.PENDING, vest_date__lte=as_of))
    for event in pending_events:
        event.status = VestingEventStatus.VESTED
        event.save(update_fields=['status', 'updated_at'])
        notify_vesting_milestone(event)
    grant.last_vesting_calculated_at = timezone.now()
    grant.save(update_fields=['last_vesting_calculated_at', 'updated_at'])


def calculate_grant_summary(grant: EquityGrant, as_of: date | None = None) -> dict:
    as_of = ensure_date(as_of or timezone.now().date())
    mark_vesting_events_as_of(grant, as_of)

    vested_units = (
        grant.vesting_events.filter(vest_date__lte=as_of, status__in=[VestingEventStatus.VESTED, VestingEventStatus.EXERCISED])
        .aggregate(total=Sum('units'))
        .get('total')
        or 0
    )
    exercised_units = grant.exercise_requests.filter(status=ExerciseRequestStatus.COMPLETED).aggregate(total=Sum('approved_units')).get('total') or 0
    forfeited_units = grant.vesting_events.filter(status=VestingEventStatus.FORFEITED).aggregate(total=Sum('units')).get('total') or grant.forfeited_units
    available_to_exercise = max(int(vested_units) - int(exercised_units), 0)
    unvested_units = max(int(grant.total_units) - int(vested_units) - int(forfeited_units), 0)

    lifecycle_status = grant.lifecycle_status
    if unvested_units == 0 and available_to_exercise == 0 and exercised_units > 0:
        lifecycle_status = GrantLifecycleStatus.EXERCISED
    elif unvested_units == 0:
        lifecycle_status = GrantLifecycleStatus.FULLY_VESTED
    elif vested_units > 0:
        lifecycle_status = GrantLifecycleStatus.ACTIVE
    grant.lifecycle_status = lifecycle_status
    grant.exercised_units = int(exercised_units)
    grant.forfeited_units = int(forfeited_units)
    grant.save(update_fields=['lifecycle_status', 'exercised_units', 'forfeited_units', 'updated_at'])

    return {
        'grant_id': str(grant.id),
        'total_units': int(grant.total_units),
        'vested_units': int(vested_units),
        'exercised_units': int(exercised_units),
        'forfeited_units': int(forfeited_units),
        'available_to_exercise': int(available_to_exercise),
        'unvested_units': int(unvested_units),
        'lifecycle_status': lifecycle_status,
    }


def run_vesting_notification_sweep(as_of: date | None = None, reminder_days: tuple[int, ...] = (30, 7, 1)) -> dict:
    as_of = ensure_date(as_of or timezone.now().date())
    matured_grants = 0
    reminder_count = 0

    grants = EquityGrant.objects.select_related('workspace', 'shareholder', 'employee__user', 'share_class').all()
    for grant in grants:
        due_events = grant.vesting_events.filter(status=VestingEventStatus.PENDING, vest_date__lte=as_of)
        if due_events.exists():
            mark_vesting_events_as_of(grant, as_of)
            matured_grants += 1

        for days_until in reminder_days:
            target_date = as_of + timedelta(days=days_until)
            upcoming_events = grant.vesting_events.filter(
                status=VestingEventStatus.PENDING,
                vest_date=target_date,
            )
            for event in upcoming_events:
                event_name = f'vesting_reminder_{days_until}d'
                already_sent = event.delivery_logs.filter(event_name=event_name, status='sent').exists()
                if already_sent:
                    continue
                notify_upcoming_vesting(event, days_until)
                reminder_count += 1

    return {
        'as_of': as_of.isoformat(),
        'grants_matured': matured_grants,
        'reminders_sent': reminder_count,
    }


def _resolve_tax_rate(grant: EquityGrant) -> Decimal:
    profile = None
    if grant.employee and hasattr(grant.employee.user, 'profile'):
        profile = grant.employee.user.profile
    elif hasattr(grant.shareholder.created_by, 'profile'):
        profile = grant.shareholder.created_by.profile
    if isinstance(profile, UserProfile):
        return Decimal(profile.tax_rate or 0)
    return Decimal('0')


@transaction.atomic
def create_exercise_request(
    grant: EquityGrant,
    requested_units: int,
    payment_method: str,
    created_by: User | None,
    notes: str = '',
) -> EquityExerciseRequest:
    summary = calculate_grant_summary(grant)
    if requested_units <= 0:
        raise ValueError('Requested units must be greater than zero.')
    if requested_units > summary['available_to_exercise']:
        raise ValueError('Requested units exceed available vested units.')

    latest_valuation = grant.workspace.equity_valuations.order_by('-valuation_date', '-created_at').first()
    market_price = Decimal(latest_valuation.price_per_share if latest_valuation else grant.exercise_price)
    strike_payment_amount = quantize_money(Decimal(requested_units) * Decimal(grant.exercise_price))
    taxable_gain = max(market_price - Decimal(grant.exercise_price), Decimal('0')) * Decimal(requested_units)
    tax_rate = _resolve_tax_rate(grant)
    withholding = quantize_money(taxable_gain * _decimal_percent(tax_rate))

    exercise_request = EquityExerciseRequest.objects.create(
        workspace=grant.workspace,
        grant=grant,
        shareholder=grant.shareholder,
        requested_units=requested_units,
        strike_price_per_unit=grant.exercise_price,
        strike_payment_amount=strike_payment_amount,
        tax_withholding_amount=withholding,
        payment_method=payment_method,
        payment_status=PaymentStatus.PROCESSING if payment_method == ExercisePaymentMethod.PAYROLL_DEDUCTION else PaymentStatus.PENDING,
        status=ExerciseRequestStatus.REQUESTED,
        due_date=(timezone.now() + timedelta(days=grant.post_termination_exercise_days)).date(),
        notes=notes,
        created_by=created_by,
    )

    approvers = []
    if grant.employee and grant.employee.manager and grant.employee.manager.user_id:
        approvers.append(grant.employee.manager.user)
    org_owner = grant.workspace.organization.owner
    if org_owner and all(candidate.id != org_owner.id for candidate in approvers):
        approvers.append(org_owner)

    if not approvers and created_by:
        approvers.append(created_by)

    for order, approver in enumerate(approvers, start=1):
        EquityExerciseApproval.objects.create(
            exercise_request=exercise_request,
            approver=approver,
            approval_order=order,
        )

    if exercise_request.approvals.count() > 1:
        exercise_request.status = ExerciseRequestStatus.FINANCE_REVIEW
        exercise_request.save(update_fields=['status', 'updated_at'])

    notify_exercise_submitted(exercise_request)

    return exercise_request


@transaction.atomic
def approve_exercise_request(exercise_request: EquityExerciseRequest, approver: User, comments: str = '') -> EquityExerciseRequest:
    approval = exercise_request.approvals.filter(approver=approver).order_by('approval_order').first()
    if approval is None:
        raise ValueError('No approval step is assigned to this approver.')
    if approval.status == ApprovalDecisionStatus.APPROVED:
        return exercise_request

    approval.status = ApprovalDecisionStatus.APPROVED
    approval.decided_at = timezone.now()
    approval.comments = comments
    approval.save(update_fields=['status', 'decided_at', 'comments', 'updated_at'])

    if exercise_request.approvals.filter(status=ApprovalDecisionStatus.REJECTED).exists():
        exercise_request.status = ExerciseRequestStatus.REJECTED
    elif exercise_request.approvals.exclude(status=ApprovalDecisionStatus.APPROVED).exists():
        exercise_request.status = ExerciseRequestStatus.LEGAL_REVIEW
    else:
        exercise_request.status = ExerciseRequestStatus.APPROVED
        exercise_request.approved_units = exercise_request.requested_units
    exercise_request.save(update_fields=['status', 'approved_units', 'updated_at'])
    if exercise_request.status == ExerciseRequestStatus.APPROVED:
        notify_exercise_status(exercise_request, approved=True)
    return exercise_request


@transaction.atomic
def reject_exercise_request(exercise_request: EquityExerciseRequest, approver: User, comments: str = '') -> EquityExerciseRequest:
    approval = exercise_request.approvals.filter(approver=approver).order_by('approval_order').first()
    if approval is None:
        raise ValueError('No approval step is assigned to this approver.')
    approval.status = ApprovalDecisionStatus.REJECTED
    approval.decided_at = timezone.now()
    approval.comments = comments
    approval.save(update_fields=['status', 'decided_at', 'comments', 'updated_at'])
    exercise_request.status = ExerciseRequestStatus.REJECTED
    exercise_request.save(update_fields=['status', 'updated_at'])
    notify_exercise_status(exercise_request, approved=False)
    return exercise_request


def _get_or_create_account(entity, code: str, name: str, account_type: str) -> ChartOfAccounts:
    account, _ = ChartOfAccounts.objects.get_or_create(
        entity=entity,
        account_code=code,
        defaults={
            'account_name': name,
            'account_type': account_type,
            'currency': entity.local_currency,
            'status': 'active',
        },
    )
    return account


@transaction.atomic
def complete_exercise_request(exercise_request: EquityExerciseRequest, actor: User | None) -> EquityExerciseRequest:
    if exercise_request.status != ExerciseRequestStatus.APPROVED:
        raise ValueError('Exercise request must be approved before completion.')

    grant = exercise_request.grant
    approved_units = exercise_request.approved_units or exercise_request.requested_units
    if approved_units <= 0:
        raise ValueError('Approved units must be greater than zero.')

    withholding_tax = quantize_money(exercise_request.tax_withholding_amount)
    tax_rate = _resolve_tax_rate(grant)
    tax_calculation = TaxCalculation.objects.create(
        entity=grant.workspace,
        tax_year=timezone.now().year,
        calculation_type='withholding',
        jurisdiction=grant.workspace.country,
        taxable_income=quantize_money(exercise_request.strike_payment_amount + withholding_tax),
        tax_rate=Decimal(tax_rate) / Decimal('100'),
        deductions={},
        credits={},
        calculated_tax=withholding_tax,
        effective_rate=Decimal(tax_rate) / Decimal('100'),
        breakdown={
            'exercise_request_id': str(exercise_request.id),
            'grant_id': str(grant.id),
            'requested_units': approved_units,
        },
    )

    cash_account = _get_or_create_account(grant.workspace, '1010', 'Equity Exercise Cash', 'asset')
    equity_account = _get_or_create_account(grant.workspace, '3000', 'Share Capital', 'equity')
    tax_account = _get_or_create_account(grant.workspace, '2100', 'Tax Withholding Payable', 'liability')
    compensation_account = _get_or_create_account(grant.workspace, '5100', 'Equity Compensation Expense', 'expense')

    reference_number = f'EX-{exercise_request.id.hex[:8].upper()}'
    journal_entry = JournalEntry.objects.create(
        entity=grant.workspace,
        entry_type='automated',
        reference_number=reference_number,
        description=f'Equity exercise for {grant.shareholder.name}',
        posting_date=timezone.now().date(),
        memo=f'Grant {grant.grant_number}',
        status='posted',
        created_by=actor,
        approved_by=actor,
        approved_at=timezone.now(),
    )

    GeneralLedger.objects.create(
        entity=grant.workspace,
        debit_account=cash_account,
        credit_account=equity_account,
        debit_amount=quantize_money(exercise_request.strike_payment_amount),
        credit_amount=quantize_money(exercise_request.strike_payment_amount),
        description='Strike price payment on option exercise',
        reference_number=reference_number,
        posting_date=timezone.now().date(),
        journal_entry=journal_entry,
        posting_status='posted',
    )

    if withholding_tax > 0:
        GeneralLedger.objects.create(
            entity=grant.workspace,
            debit_account=compensation_account,
            credit_account=tax_account,
            debit_amount=withholding_tax,
            credit_amount=withholding_tax,
            description='Tax withholding on equity exercise',
            reference_number=reference_number,
            posting_date=timezone.now().date(),
            journal_entry=journal_entry,
            posting_status='posted',
        )

    holding, _ = EquityHolding.objects.get_or_create(
        workspace=grant.workspace,
        shareholder=grant.shareholder,
        share_class=grant.share_class,
        defaults={
            'quantity': 0,
            'diluted_quantity': 0,
            'ownership_percent': 0,
            'issued_at': timezone.now().date(),
            'strike_price': grant.exercise_price,
        },
    )
    holding.quantity += approved_units
    holding.diluted_quantity += approved_units
    holding.issued_at = holding.issued_at or timezone.now().date()
    holding.strike_price = grant.exercise_price
    holding.save()

    grant.share_class.issued_shares += approved_units
    grant.share_class.save(update_fields=['issued_shares', 'updated_at'])

    EquityTransaction.objects.create(
        workspace=grant.workspace,
        transaction_type=TransactionType.EXERCISE,
        shareholder=grant.shareholder,
        share_class=grant.share_class,
        quantity=approved_units,
        price_per_share=grant.exercise_price,
        effective_date=timezone.now().date(),
        approval_status='executed',
        compliance_checked=True,
        digital_signature_required=False,
        audit_metadata={
            'exercise_request_id': str(exercise_request.id),
            'grant_id': str(grant.id),
            'journal_entry_id': journal_entry.id,
            'tax_calculation_id': tax_calculation.id,
        },
        created_by=actor,
    )

    certificate = EquityShareCertificate.objects.create(
        workspace=grant.workspace,
        exercise_request=exercise_request,
        grant=grant,
        certificate_number=f'CERT-{grant.grant_number}-{timezone.now().strftime("%Y%m%d%H%M%S")}',
        issued_to=grant.shareholder,
        share_class=grant.share_class,
        issued_units=approved_units,
        issue_date=timezone.now().date(),
        status=CertificateStatus.ISSUED,
        certificate_payload={
            'shareholder': grant.shareholder.name,
            'grant_number': grant.grant_number,
            'units': approved_units,
            'issue_date': timezone.now().date().isoformat(),
        },
        issued_by=actor,
    )

    payroll_tax_event = EquityPayrollTaxEvent.objects.create(
        workspace=grant.workspace,
        grant=grant,
        exercise_request=exercise_request,
        staff=grant.employee,
        event_type='exercise_withholding',
        gross_amount=quantize_money(exercise_request.strike_payment_amount),
        withholding_amount=withholding_tax,
        payroll_sync_status='queued' if exercise_request.payment_method == ExercisePaymentMethod.PAYROLL_DEDUCTION else 'synced',
        tax_jurisdiction=grant.workspace.country,
        reference_number=reference_number,
        source_account=compensation_account,
        destination_account=tax_account,
        details={
            'certificate_id': str(certificate.id),
            'payment_method': exercise_request.payment_method,
        },
    )

    payment_sync_result = sync_exercise_payment(exercise_request)
    payroll_sync_result = sync_payroll_tax_event(payroll_tax_event)

    exercise_request.tax_calculation = tax_calculation
    exercise_request.journal_entry = journal_entry
    exercise_request.status = ExerciseRequestStatus.COMPLETED
    exercise_request.payment_status = PaymentStatus.PAID if exercise_request.payment_status != PaymentStatus.WAIVED else PaymentStatus.WAIVED
    exercise_request.exercise_date = timezone.now().date()
    exercise_request.completed_at = timezone.now()
    exercise_request.approved_units = approved_units
    exercise_request.save(
        update_fields=[
            'tax_calculation',
            'journal_entry',
            'status',
            'payment_status',
            'exercise_date',
            'completed_at',
            'approved_units',
            'updated_at',
        ]
    )

    exercised_events = grant.vesting_events.filter(status=VestingEventStatus.VESTED).order_by('vest_date', 'created_at')
    units_remaining = approved_units
    for event in exercised_events:
        if units_remaining <= 0:
            break
        if event.units <= units_remaining:
            event.status = VestingEventStatus.EXERCISED
            units_remaining -= event.units
            event.save(update_fields=['status', 'updated_at'])

    if payment_sync_result.get('results') or payroll_sync_result.get('results'):
        exercise_request.notes = '\n'.join(
            item for item in [
                exercise_request.notes.strip(),
                f"Payment sync: {'ok' if payment_sync_result.get('successful') else 'pending/failed'} ({payment_sync_result.get('configured_adapters', 0)} adapters)",
                f"Payroll sync: {'ok' if payroll_sync_result.get('successful') else 'pending/failed'} ({payroll_sync_result.get('configured_adapters', 0)} adapters)",
            ]
            if item
        )
        exercise_request.save(update_fields=['notes', 'updated_at'])

    notify_certificate_released(certificate)
    calculate_grant_summary(grant)
    return exercise_request
