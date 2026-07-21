from __future__ import annotations

from datetime import timedelta
from decimal import Decimal, ROUND_DOWN

from django.db import transaction
from django.db.models import Q
from django.utils import timezone

from atonixcorp.platform_foundation import log_platform_audit_event, sync_equity_scenario_approval_task

from .automation import (
    notify_scenario_escalated,
    notify_scenario_reminder,
    notify_scenario_committed,
    notify_scenario_requested,
    notify_scenario_review_decision,
)
from .models import (
    AntiDilutionType,
    ApprovalStatus,
    EquityFundingRound,
    EquityGrant,
    EquityHolding,
    EquityOptionPoolReserve,
    EquityReport,
    EquityScenarioApprovalEvent,
    EquityScenarioApprovalPolicy,
    EquityScenarioApproval,
    EquityShareClass,
    EquityShareholder,
    EquityTransaction,
    ScenarioApprovalEventType,
    ScenarioApprovalStatus,
    ScenarioReviewStatus,
    ShareClassType,
    ShareholderType,
    TransactionType,
)


MONEY_QUANTUM = Decimal('0.01')
SHARE_QUANTUM = Decimal('1')


def _to_decimal(value, default: str = '0') -> Decimal:
    if value in (None, ''):
        return Decimal(default)
    if isinstance(value, Decimal):
        return value
    return Decimal(str(value))


def _money(value) -> Decimal:
    return _to_decimal(value).quantize(MONEY_QUANTUM)


def _round_shares(value) -> int:
    decimal_value = _to_decimal(value)
    if decimal_value <= 0:
        return 0
    return int(decimal_value.quantize(SHARE_QUANTUM, rounding=ROUND_DOWN))


def _serialize_money(value) -> str:
    return str(_money(value))


def _serialize_ratio(value) -> str:
    return str(_to_decimal(value).quantize(Decimal('0.0001')))


def get_or_create_scenario_approval_policy(workspace):
    policy, _ = EquityScenarioApprovalPolicy.objects.get_or_create(workspace=workspace)
    return policy


def _log_approval_event(approval, event_type, title, message='', *, actor=None, metadata=None):
    return EquityScenarioApprovalEvent.objects.create(
        approval=approval,
        actor=actor,
        event_type=event_type,
        title=title,
        message=message,
        metadata=metadata or {},
    )


def _validate_submission_policy(policy):
    errors = {}
    if policy.require_explicit_reviewers:
        if not policy.board_reviewers.exists():
            errors['board_reviewers'] = 'Board reviewers must be assigned before a scenario can be submitted.'
        if not policy.legal_reviewers.exists():
            errors['legal_reviewers'] = 'Legal reviewers must be assigned before a scenario can be submitted.'
    if policy.require_designated_backups:
        if not policy.board_escalation_reviewers.exists():
            errors['board_escalation_reviewers'] = 'Board backup reviewers must be assigned before a scenario can be submitted.'
        if not policy.legal_escalation_reviewers.exists():
            errors['legal_escalation_reviewers'] = 'Legal backup reviewers must be assigned before a scenario can be submitted.'
    if errors:
        raise ValueError(errors)


def _workspace_staff_for_user(workspace, user):
    if not user or not getattr(user, 'is_authenticated', False):
        return None
    return workspace.staff.select_related('role', 'department').filter(user=user, status='active').first()


def _policy_reviewer_users(staff_queryset, fallback_user=None):
    users = []
    seen = set()
    for staff_member in staff_queryset.select_related('user'):
        if not staff_member.user_id or staff_member.user_id in seen:
            continue
        seen.add(staff_member.user_id)
        users.append(staff_member.user)
    if not users and fallback_user:
        users.append(fallback_user)
    return users


def get_board_review_users(workspace):
    policy = get_or_create_scenario_approval_policy(workspace)
    return _policy_reviewer_users(policy.board_reviewers.all(), fallback_user=getattr(workspace.organization, 'owner', None))


def get_legal_review_users(workspace):
    policy = get_or_create_scenario_approval_policy(workspace)
    return _policy_reviewer_users(policy.legal_reviewers.all(), fallback_user=getattr(workspace.organization, 'owner', None))


def get_board_escalation_users(workspace):
    policy = get_or_create_scenario_approval_policy(workspace)
    escalation_queryset = policy.board_escalation_reviewers.all()
    if escalation_queryset.exists():
        return _policy_reviewer_users(escalation_queryset, fallback_user=getattr(workspace.organization, 'owner', None))
    return get_board_review_users(workspace)


def get_legal_escalation_users(workspace):
    policy = get_or_create_scenario_approval_policy(workspace)
    escalation_queryset = policy.legal_escalation_reviewers.all()
    if escalation_queryset.exists():
        return _policy_reviewer_users(escalation_queryset, fallback_user=getattr(workspace.organization, 'owner', None))
    return get_legal_review_users(workspace)


def user_can_board_approve(workspace, user) -> bool:
    if getattr(workspace.organization, 'owner_id', None) == getattr(user, 'id', None):
        return True
    staff_member = _workspace_staff_for_user(workspace, user)
    if not staff_member:
        return False
    policy = get_or_create_scenario_approval_policy(workspace)
    return policy.board_reviewers.filter(pk=staff_member.pk).exists()


def user_can_legal_approve(workspace, user) -> bool:
    if getattr(workspace.organization, 'owner_id', None) == getattr(user, 'id', None):
        return True
    staff_member = _workspace_staff_for_user(workspace, user)
    if not staff_member:
        return False
    policy = get_or_create_scenario_approval_policy(workspace)
    return policy.legal_reviewers.filter(pk=staff_member.pk).exists()


def _apply_policy_sla_dates(approval: EquityScenarioApproval, policy: EquityScenarioApprovalPolicy):
    base_time = approval.created_at or timezone.now()
    approval.board_due_at = base_time + timedelta(hours=int(policy.board_sla_hours or 72))
    approval.legal_due_at = base_time + timedelta(hours=int(policy.legal_sla_hours or 72))


def get_scenario_approval_inbox(workspace, user):
    queryset = EquityScenarioApproval.objects.filter(workspace=workspace).select_related('requested_by', 'board_approved_by', 'legal_approved_by', 'committed_round').prefetch_related('events__actor')
    is_board_reviewer = user_can_board_approve(workspace, user)
    is_legal_reviewer = user_can_legal_approve(workspace, user)
    is_owner = getattr(workspace.organization, 'owner_id', None) == getattr(user, 'id', None)

    if is_owner:
        scoped_queryset = queryset.exclude(status=ScenarioApprovalStatus.COMMITTED)
    else:
        filters = []
        if is_board_reviewer:
            filters.append({'board_status': ScenarioReviewStatus.PENDING})
        if is_legal_reviewer:
            filters.append({'legal_status': ScenarioReviewStatus.PENDING})
        filters.append({'requested_by': user})
        filters.append({'status': ScenarioApprovalStatus.APPROVED})
        query = Q()
        for item in filters:
            query |= Q(**item)
        scoped_queryset = queryset.filter(query).exclude(status=ScenarioApprovalStatus.COMMITTED).distinct()

    now = timezone.now()
    pending = list(scoped_queryset.order_by('-created_at'))
    overdue = [
        item for item in pending
        if (item.board_status == ScenarioReviewStatus.PENDING and item.board_due_at and item.board_due_at <= now)
        or (item.legal_status == ScenarioReviewStatus.PENDING and item.legal_due_at and item.legal_due_at <= now)
    ]
    return {
        'pending': pending,
        'overdue': overdue,
        'summary': {
            'pending_count': len(pending),
            'overdue_count': len(overdue),
            'board_reviewer': is_board_reviewer,
            'legal_reviewer': is_legal_reviewer,
            'owner_reviewer': is_owner,
        },
    }


def run_scenario_approval_sla_sweep(*, workspace=None, as_of=None):
    as_of = as_of or timezone.now()
    approvals = EquityScenarioApproval.objects.select_related('workspace', 'requested_by', 'board_approved_by', 'legal_approved_by').all()
    if workspace is not None:
        approvals = approvals.filter(workspace=workspace)

    reminder_count = 0
    escalation_count = 0
    for approval in approvals:
        policy = get_or_create_scenario_approval_policy(approval.workspace)
        update_fields = []
        if approval.board_due_at is None or approval.legal_due_at is None:
            _apply_policy_sla_dates(approval, policy)
            update_fields.extend(['board_due_at', 'legal_due_at'])

        reminder_interval = timedelta(hours=int(policy.reminder_frequency_hours or 24))
        escalation_grace = timedelta(hours=int(policy.escalation_grace_hours or 24))

        if approval.board_status == ScenarioReviewStatus.PENDING and approval.board_due_at:
            should_remind = approval.board_due_at <= as_of and (not approval.board_last_reminder_at or approval.board_last_reminder_at + reminder_interval <= as_of)
            if should_remind:
                notify_scenario_reminder(approval, reviewer_type='board', escalated=False)
                _log_approval_event(
                    approval,
                    ScenarioApprovalEventType.REMINDER,
                    'Board review reminder sent',
                    'Board review is past due and a reminder was sent.',
                    metadata={'reviewer_type': 'board'},
                )
                approval.board_last_reminder_at = as_of
                update_fields.append('board_last_reminder_at')
                reminder_count += 1
            should_escalate = policy.escalation_enabled and approval.board_due_at + escalation_grace <= as_of and approval.board_escalated_at is None
            if should_escalate:
                notify_scenario_escalated(approval, reviewer_type='board')
                _log_approval_event(
                    approval,
                    ScenarioApprovalEventType.ESCALATED,
                    'Board review escalated',
                    'Board review remained overdue and was escalated to backup reviewers.',
                    metadata={'reviewer_type': 'board'},
                )
                approval.board_escalated_at = as_of
                update_fields.append('board_escalated_at')
                escalation_count += 1

        if approval.legal_status == ScenarioReviewStatus.PENDING and approval.legal_due_at:
            should_remind = approval.legal_due_at <= as_of and (not approval.legal_last_reminder_at or approval.legal_last_reminder_at + reminder_interval <= as_of)
            if should_remind:
                notify_scenario_reminder(approval, reviewer_type='legal', escalated=False)
                _log_approval_event(
                    approval,
                    ScenarioApprovalEventType.REMINDER,
                    'Legal review reminder sent',
                    'Legal review is past due and a reminder was sent.',
                    metadata={'reviewer_type': 'legal'},
                )
                approval.legal_last_reminder_at = as_of
                update_fields.append('legal_last_reminder_at')
                reminder_count += 1
            should_escalate = policy.escalation_enabled and approval.legal_due_at + escalation_grace <= as_of and approval.legal_escalated_at is None
            if should_escalate:
                notify_scenario_escalated(approval, reviewer_type='legal')
                _log_approval_event(
                    approval,
                    ScenarioApprovalEventType.ESCALATED,
                    'Legal review escalated',
                    'Legal review remained overdue and was escalated to backup reviewers.',
                    metadata={'reviewer_type': 'legal'},
                )
                approval.legal_escalated_at = as_of
                update_fields.append('legal_escalated_at')
                escalation_count += 1

        if update_fields:
            approval.save(update_fields=[*set(update_fields), 'updated_at'])

    return {
        'processed': approvals.count(),
        'reminders_sent': reminder_count,
        'escalations_sent': escalation_count,
    }


def _allocate_proportional_shares(total_shares: int, weights: list[Decimal]) -> list[int]:
    if total_shares <= 0 or not weights:
        return [0 for _ in weights]
    weight_sum = sum(weights)
    if weight_sum <= 0:
        return [0 for _ in weights]

    raw_allocations = [Decimal(total_shares) * weight / weight_sum for weight in weights]
    base_allocations = [int(value.quantize(SHARE_QUANTUM, rounding=ROUND_DOWN)) for value in raw_allocations]
    remainder = total_shares - sum(base_allocations)
    ranked_indices = sorted(
        range(len(weights)),
        key=lambda index: (raw_allocations[index] - Decimal(base_allocations[index]), weights[index]),
        reverse=True,
    )
    for index in ranked_indices[:remainder]:
        base_allocations[index] += 1
    return base_allocations


def build_cap_table_snapshot(workspace) -> dict:
    holdings = list(
        EquityHolding.objects.filter(workspace=workspace)
        .select_related('shareholder', 'share_class')
        .order_by('share_class__liquidation_seniority', 'share_class__name', 'shareholder__name')
    )
    outstanding_grants = list(
        EquityGrant.objects.filter(workspace=workspace)
        .select_related('share_class')
        .order_by('share_class__name', 'grant_number')
    )
    option_pool_reserves = list(
        EquityOptionPoolReserve.objects.filter(workspace=workspace)
        .select_related('share_class', 'funding_round')
        .order_by('share_class__name', '-updated_at')
    )

    share_classes = {str(item.id): item for item in EquityShareClass.objects.filter(workspace=workspace)}
    transaction_prices = {}
    for transaction in EquityTransaction.objects.filter(workspace=workspace, shareholder__isnull=False, share_class__isnull=False).order_by('-effective_date', '-created_at'):
        transaction_prices.setdefault((str(transaction.shareholder_id), str(transaction.share_class_id)), _to_decimal(transaction.price_per_share))
    class_rows = {}
    holder_rows = []
    total_outstanding = 0
    total_fully_diluted = 0
    total_invested = Decimal('0')

    for holding in holdings:
        share_count = int(holding.quantity or 0)
        diluted_count = int(holding.diluted_quantity or 0) or share_count
        fallback_price = transaction_prices.get((str(holding.shareholder_id), str(holding.share_class_id)), Decimal('0'))
        issue_price = _to_decimal(holding.issue_price_per_share or holding.strike_price or fallback_price)
        invested_amount = _money(holding.invested_amount or (issue_price * Decimal(share_count)))
        total_outstanding += share_count
        total_fully_diluted += diluted_count
        total_invested += invested_amount

        class_key = str(holding.share_class_id)
        class_row = class_rows.setdefault(
            class_key,
            {
                'id': class_key,
                'name': holding.share_class.name,
                'class_type': holding.share_class.class_type,
                'outstanding_shares': 0,
                'fully_diluted_shares': 0,
                'invested_capital': Decimal('0'),
                'preference_multiple': _to_decimal(holding.share_class.preference_multiple),
                'participating_preference': holding.share_class.participating_preference,
                'participation_cap_multiple': _to_decimal(holding.share_class.participation_cap_multiple),
                'liquidation_seniority': int(holding.share_class.liquidation_seniority or 0),
                'conversion_price': _to_decimal(holding.share_class.conversion_price),
                'anti_dilution_type': holding.share_class.anti_dilution_type,
                'anti_dilution_basis': holding.share_class.anti_dilution_basis,
                'pro_rata_rights': holding.share_class.pro_rata_rights,
            },
        )
        class_row['outstanding_shares'] += share_count
        class_row['fully_diluted_shares'] += diluted_count
        class_row['invested_capital'] += invested_amount

        holder_rows.append(
            {
                'holder_id': str(holding.shareholder_id),
                'holder_name': holding.shareholder.name,
                'share_class_id': class_key,
                'share_class_name': holding.share_class.name,
                'shares': share_count,
                'fully_diluted_shares': diluted_count,
                'invested_capital': invested_amount,
                'issue_price_per_share': issue_price,
                'pro_rata_eligible': bool(holding.pro_rata_eligible or holding.share_class.pro_rata_rights),
                'pro_rata_take_up_percent': _to_decimal(holding.pro_rata_take_up_percent or 100),
                'scenario_source': 'holding',
            }
        )

    grant_pool_by_class = {}
    for grant in outstanding_grants:
        unissued_units = max(int(grant.total_units or 0) - int(grant.exercised_units or 0) - int(grant.forfeited_units or 0), 0)
        if not unissued_units:
            continue
        class_key = str(grant.share_class_id)
        grant_pool_by_class[class_key] = grant_pool_by_class.get(class_key, 0) + unissued_units

    for class_key, grant_units in grant_pool_by_class.items():
        share_class = share_classes.get(class_key)
        if share_class is None:
            continue
        total_fully_diluted += grant_units
        class_row = class_rows.setdefault(
            class_key,
            {
                'id': class_key,
                'name': share_class.name,
                'class_type': share_class.class_type,
                'outstanding_shares': 0,
                'fully_diluted_shares': 0,
                'invested_capital': Decimal('0'),
                'preference_multiple': _to_decimal(share_class.preference_multiple),
                'participating_preference': share_class.participating_preference,
                'participation_cap_multiple': _to_decimal(share_class.participation_cap_multiple),
                'liquidation_seniority': int(share_class.liquidation_seniority or 0),
                'conversion_price': _to_decimal(share_class.conversion_price),
                'anti_dilution_type': share_class.anti_dilution_type,
                'anti_dilution_basis': share_class.anti_dilution_basis,
                'pro_rata_rights': share_class.pro_rata_rights,
            },
        )
        class_row['fully_diluted_shares'] += grant_units
        holder_rows.append(
            {
                'holder_id': f'grant-pool-{class_key}',
                'holder_name': f'{share_class.name} Grant Pool',
                'share_class_id': class_key,
                'share_class_name': share_class.name,
                'shares': 0,
                'fully_diluted_shares': grant_units,
                'invested_capital': Decimal('0'),
                'issue_price_per_share': Decimal('0'),
                'pro_rata_eligible': False,
                'pro_rata_take_up_percent': Decimal('0'),
                'scenario_source': 'grant_pool',
            }
        )

    for reserve in option_pool_reserves:
        available_units = max(int(reserve.reserved_shares or 0) - int(reserve.allocated_shares or 0), 0)
        if not available_units:
            continue
        class_key = str(reserve.share_class_id)
        total_fully_diluted += available_units
        class_row = class_rows.setdefault(
            class_key,
            {
                'id': class_key,
                'name': reserve.share_class.name,
                'class_type': reserve.share_class.class_type,
                'outstanding_shares': 0,
                'fully_diluted_shares': 0,
                'invested_capital': Decimal('0'),
                'preference_multiple': _to_decimal(reserve.share_class.preference_multiple),
                'participating_preference': reserve.share_class.participating_preference,
                'participation_cap_multiple': _to_decimal(reserve.share_class.participation_cap_multiple),
                'liquidation_seniority': int(reserve.share_class.liquidation_seniority or 0),
                'conversion_price': _to_decimal(reserve.share_class.conversion_price),
                'anti_dilution_type': reserve.share_class.anti_dilution_type,
                'anti_dilution_basis': reserve.share_class.anti_dilution_basis,
                'pro_rata_rights': reserve.share_class.pro_rata_rights,
            },
        )
        class_row['fully_diluted_shares'] += available_units
        holder_rows.append(
            {
                'holder_id': f'option-pool-reserve-{reserve.id}',
                'holder_name': f'{reserve.share_class.name} Option Reserve',
                'share_class_id': class_key,
                'share_class_name': reserve.share_class.name,
                'shares': 0,
                'fully_diluted_shares': available_units,
                'invested_capital': Decimal('0'),
                'issue_price_per_share': Decimal('0'),
                'pro_rata_eligible': False,
                'pro_rata_take_up_percent': Decimal('0'),
                'scenario_source': 'option_pool_reserve',
            }
        )

    holder_rows.sort(key=lambda row: (row['share_class_name'], row['holder_name']))
    class_list = list(class_rows.values())
    class_list.sort(key=lambda row: (-row['liquidation_seniority'], row['name']))

    for row in holder_rows:
        base = total_fully_diluted or 1
        row['ownership_percent'] = (_to_decimal(row['fully_diluted_shares']) * Decimal('100') / Decimal(base)).quantize(Decimal('0.0001'))

    for row in class_list:
        base = total_fully_diluted or 1
        row['ownership_percent'] = (_to_decimal(row['fully_diluted_shares']) * Decimal('100') / Decimal(base)).quantize(Decimal('0.0001'))

    return {
        'outstanding_shares': total_outstanding,
        'fully_diluted_shares': total_fully_diluted,
        'invested_capital': total_invested,
        'holder_rows': holder_rows,
        'class_rows': class_list,
    }


def _scenario_security(snapshot: dict, scenario_input: dict) -> dict:
    share_class_id = str(scenario_input.get('share_class') or '')
    class_lookup = {row['id']: row for row in snapshot['class_rows']}
    selected = class_lookup.get(share_class_id)
    if selected:
        return dict(selected)
    return {
        'id': 'scenario-security',
        'name': scenario_input.get('name') or 'Scenario Security',
        'class_type': 'preferred',
        'outstanding_shares': 0,
        'fully_diluted_shares': 0,
        'invested_capital': Decimal('0'),
        'preference_multiple': Decimal('1'),
        'participating_preference': False,
        'participation_cap_multiple': Decimal('0'),
        'liquidation_seniority': 1,
        'conversion_price': Decimal('0'),
        'anti_dilution_type': AntiDilutionType.NONE,
        'anti_dilution_basis': 'broad_based',
        'pro_rata_rights': True,
    }


def _apply_anti_dilution(snapshot: dict, price_per_share: Decimal, amount_raised: Decimal, new_money_shares: int) -> tuple[list[dict], dict[str, int]]:
    adjustments = []
    incremental_shares_by_class = {}
    fully_diluted_before = snapshot['fully_diluted_shares']
    outstanding_before = snapshot['outstanding_shares']

    class_lookup = {row['id']: row for row in snapshot['class_rows']}
    holders_by_class = {}
    for row in snapshot['holder_rows']:
        holders_by_class.setdefault(row['share_class_id'], []).append(row)

    for class_row in snapshot['class_rows']:
        anti_dilution_type = class_row['anti_dilution_type']
        old_conversion_price = _to_decimal(class_row['conversion_price'])
        if anti_dilution_type == AntiDilutionType.NONE or old_conversion_price <= 0 or price_per_share >= old_conversion_price:
            continue

        basis_shares = fully_diluted_before if class_row['anti_dilution_basis'] == 'broad_based' else outstanding_before
        if anti_dilution_type == AntiDilutionType.FULL_RATCHET:
            new_conversion_price = price_per_share
        else:
            basis = Decimal(basis_shares or 1)
            equivalent_old_shares = amount_raised / old_conversion_price if old_conversion_price else Decimal('0')
            new_conversion_price = old_conversion_price * ((basis + equivalent_old_shares) / (basis + Decimal(new_money_shares or 0)))

        if new_conversion_price <= 0 or new_conversion_price >= old_conversion_price:
            continue

        new_ratio = old_conversion_price / new_conversion_price
        class_shares = int(class_row['fully_diluted_shares'] or 0)
        incremental_shares = _round_shares((Decimal(class_shares) * new_ratio) - Decimal(class_shares))
        if incremental_shares <= 0:
            continue

        incremental_shares_by_class[class_row['id']] = incremental_shares
        holder_weights = [Decimal(holder['fully_diluted_shares']) for holder in holders_by_class.get(class_row['id'], [])]
        holder_allocations = _allocate_proportional_shares(incremental_shares, holder_weights)
        holder_adjustments = []
        for holder, allocation in zip(holders_by_class.get(class_row['id'], []), holder_allocations):
            if allocation <= 0:
                continue
            holder_adjustments.append({'holder_id': holder['holder_id'], 'holder_name': holder['holder_name'], 'added_shares': allocation})

        adjustments.append(
            {
                'share_class_id': class_row['id'],
                'share_class_name': class_row['name'],
                'anti_dilution_type': anti_dilution_type,
                'old_conversion_price': _serialize_ratio(old_conversion_price),
                'new_conversion_price': _serialize_ratio(new_conversion_price),
                'incremental_shares': incremental_shares,
                'holder_adjustments': holder_adjustments,
            }
        )

    return adjustments, incremental_shares_by_class


def _build_post_financing_positions(snapshot: dict, scenario_input: dict) -> tuple[list[dict], dict]:
    pre_money = _money(scenario_input['pre_money_valuation'])
    amount_raised = _money(scenario_input['amount_raised'])
    option_pool_top_up = int(scenario_input.get('option_pool_top_up') or 0)
    apply_pro_rata = bool(scenario_input.get('apply_pro_rata', True))
    include_anti_dilution = bool(scenario_input.get('include_anti_dilution', True))
    investor_name = scenario_input.get('investor_name') or 'New Investor'

    security = _scenario_security(snapshot, scenario_input)
    fully_diluted_before = snapshot['fully_diluted_shares'] or 0
    default_price = (pre_money / Decimal(fully_diluted_before)).quantize(Decimal('0.0001')) if fully_diluted_before else Decimal('0')
    if default_price <= 0:
        default_price = _to_decimal(security.get('conversion_price')) or Decimal('1')
    price_per_share = _to_decimal(scenario_input.get('price_per_share')) if scenario_input.get('price_per_share') not in (None, '') else default_price
    if price_per_share <= 0:
        price_per_share = Decimal('1')
    new_money_shares = _round_shares(amount_raised / price_per_share)

    anti_dilution_adjustments, anti_dilution_map = _apply_anti_dilution(snapshot, price_per_share, amount_raised, new_money_shares) if include_anti_dilution else ([], {})

    positions = [
        {
            **row,
            'post_shares': int(row['fully_diluted_shares'] or 0),
            'new_money_shares': 0,
            'pro_rata_shares': 0,
            'anti_dilution_shares': 0,
            'option_pool_shares': 0,
            'scenario_source': row['scenario_source'],
        }
        for row in snapshot['holder_rows']
    ]

    position_lookup = {(row['holder_id'], row['share_class_id']): row for row in positions}

    for class_id, added_shares in anti_dilution_map.items():
        class_positions = [row for row in positions if row['share_class_id'] == class_id]
        allocations = _allocate_proportional_shares(added_shares, [Decimal(row['post_shares']) for row in class_positions])
        for row, allocation in zip(class_positions, allocations):
            row['post_shares'] += allocation
            row['anti_dilution_shares'] += allocation

    eligible_positions = [row for row in positions if row['pro_rata_eligible']]
    pro_rata_weights = [Decimal(row['post_shares']) * (row['pro_rata_take_up_percent'] / Decimal('100')) for row in eligible_positions]
    pro_rata_allocations = _allocate_proportional_shares(new_money_shares, pro_rata_weights) if apply_pro_rata else [0 for _ in eligible_positions]

    for row, allocation in zip(eligible_positions, pro_rata_allocations):
        if allocation <= 0:
            continue
        target_key = (row['holder_id'], security['id'])
        target_position = position_lookup.get(target_key)
        if target_position is None:
            target_position = {
                'holder_id': row['holder_id'],
                'holder_name': row['holder_name'],
                'share_class_id': security['id'],
                'share_class_name': security['name'],
                'shares': 0,
                'fully_diluted_shares': 0,
                'invested_capital': Decimal('0'),
                'issue_price_per_share': price_per_share,
                'pro_rata_eligible': False,
                'pro_rata_take_up_percent': Decimal('0'),
                'post_shares': 0,
                'new_money_shares': 0,
                'pro_rata_shares': 0,
                'anti_dilution_shares': 0,
                'option_pool_shares': 0,
                'scenario_source': 'pro_rata',
            }
            positions.append(target_position)
            position_lookup[target_key] = target_position
        target_position['post_shares'] += allocation
        target_position['new_money_shares'] += allocation
        target_position['pro_rata_shares'] += allocation
        target_position['invested_capital'] += _money(price_per_share * Decimal(allocation))

    pro_rata_total = sum(pro_rata_allocations)
    investor_allocation = max(new_money_shares - pro_rata_total, 0)
    if investor_allocation:
        positions.append(
            {
                'holder_id': 'scenario-new-investor',
                'holder_name': investor_name,
                'share_class_id': security['id'],
                'share_class_name': security['name'],
                'shares': 0,
                'fully_diluted_shares': 0,
                'invested_capital': _money(price_per_share * Decimal(investor_allocation)),
                'issue_price_per_share': price_per_share,
                'pro_rata_eligible': False,
                'pro_rata_take_up_percent': Decimal('0'),
                'post_shares': investor_allocation,
                'new_money_shares': investor_allocation,
                'pro_rata_shares': 0,
                'anti_dilution_shares': 0,
                'option_pool_shares': 0,
                'scenario_source': 'new_investor',
            }
        )

    if option_pool_top_up:
        positions.append(
            {
                'holder_id': 'scenario-option-pool',
                'holder_name': 'Scenario Option Pool',
                'share_class_id': security['id'],
                'share_class_name': f"{security['name']} Option Pool",
                'shares': 0,
                'fully_diluted_shares': 0,
                'invested_capital': Decimal('0'),
                'issue_price_per_share': Decimal('0'),
                'pro_rata_eligible': False,
                'pro_rata_take_up_percent': Decimal('0'),
                'post_shares': option_pool_top_up,
                'new_money_shares': 0,
                'pro_rata_shares': 0,
                'anti_dilution_shares': 0,
                'option_pool_shares': option_pool_top_up,
                'scenario_source': 'option_pool',
            }
        )

    total_post_shares = sum(int(row['post_shares']) for row in positions)
    for row in positions:
        row['ownership_percent'] = (_to_decimal(row['post_shares']) * Decimal('100') / Decimal(total_post_shares or 1)).quantize(Decimal('0.0001'))

    return positions, {
        'pre_money_valuation': _serialize_money(pre_money),
        'post_money_valuation': _serialize_money(pre_money + amount_raised),
        'price_per_share': _serialize_ratio(price_per_share),
        'new_money_shares': new_money_shares,
        'pro_rata_shares': pro_rata_total,
        'new_investor_shares': investor_allocation,
        'option_pool_top_up': option_pool_top_up,
        'anti_dilution_adjustments': anti_dilution_adjustments,
        'security': {
            'id': security['id'],
            'name': security['name'],
            'class_type': security['class_type'],
        },
    }


def _waterfall_from_positions(positions: list[dict], snapshot: dict, exit_value: Decimal, financing_meta: dict) -> dict:
    class_lookup = {row['id']: row for row in snapshot['class_rows']}
    security_id = financing_meta['security']['id']
    if security_id not in class_lookup:
        class_lookup[security_id] = {
            'id': security_id,
            'name': financing_meta['security']['name'],
            'class_type': financing_meta['security']['class_type'],
            'preference_multiple': Decimal('1') if financing_meta['security']['class_type'] == 'preferred' else Decimal('0'),
            'participating_preference': False,
            'participation_cap_multiple': Decimal('0'),
            'liquidation_seniority': 1,
            'conversion_price': _to_decimal(financing_meta['price_per_share']),
        }

    class_rows = {}
    total_as_converted = 0
    for row in positions:
        class_term = class_lookup.get(row['share_class_id'], {})
        entry = class_rows.setdefault(
            row['share_class_id'],
            {
                'share_class_id': row['share_class_id'],
                'share_class_name': row['share_class_name'],
                'shares': 0,
                'invested_capital': Decimal('0'),
                'preference_multiple': _to_decimal(class_term.get('preference_multiple', 0)),
                'participating_preference': bool(class_term.get('participating_preference', False)),
                'participation_cap_multiple': _to_decimal(class_term.get('participation_cap_multiple', 0)),
                'liquidation_seniority': int(class_term.get('liquidation_seniority', 0) or 0),
                'residual_paid': Decimal('0'),
                'preference_paid': Decimal('0'),
            },
        )
        entry['shares'] += int(row['post_shares'])
        entry['invested_capital'] += _money(row['invested_capital'])
        total_as_converted += int(row['post_shares'])

    for class_row in class_rows.values():
        shares = Decimal(class_row['shares'] or 0)
        class_row['as_converted_value'] = _money(exit_value * shares / Decimal(total_as_converted or 1))
        class_row['preference_claim'] = _money(class_row['invested_capital'] * class_row['preference_multiple'])
        class_row['takes_preference'] = bool(class_row['participating_preference']) or class_row['preference_claim'] > class_row['as_converted_value']

    remaining_exit = _money(exit_value)
    elected_preference = [row for row in class_rows.values() if row['takes_preference'] and row['preference_claim'] > 0]
    elected_preference.sort(key=lambda row: -row['liquidation_seniority'])

    for class_row in elected_preference:
        payout = min(remaining_exit, class_row['preference_claim'])
        class_row['preference_paid'] = payout
        remaining_exit = _money(remaining_exit - payout)

    residual_candidates = [row for row in class_rows.values() if row['participating_preference'] or not row['takes_preference']]
    uncapped = {row['share_class_id'] for row in residual_candidates}
    residual_remaining = remaining_exit
    while residual_remaining > 0 and uncapped:
        pool_shares = sum(class_rows[class_id]['shares'] for class_id in uncapped)
        if pool_shares <= 0:
            break
        distributed = Decimal('0')
        capped_this_round = set()
        for class_id in list(uncapped):
            class_row = class_rows[class_id]
            allocation = _money(residual_remaining * Decimal(class_row['shares']) / Decimal(pool_shares))
            cap_multiple = class_row['participation_cap_multiple']
            if class_row['participating_preference'] and cap_multiple > 0:
                cap_total = _money(class_row['invested_capital'] * cap_multiple)
                available = max(cap_total - class_row['preference_paid'] - class_row['residual_paid'], Decimal('0'))
                if allocation >= available:
                    allocation = _money(available)
                    capped_this_round.add(class_id)
            class_row['residual_paid'] += allocation
            distributed += allocation
        residual_remaining = _money(residual_remaining - distributed)
        if not capped_this_round:
            break
        uncapped -= capped_this_round

    class_distributions = []
    holder_distributions = []
    class_totals = {key: _money(row['preference_paid'] + row['residual_paid']) for key, row in class_rows.items()}
    for class_id, class_row in sorted(class_rows.items(), key=lambda item: (-item[1]['liquidation_seniority'], item[1]['share_class_name'])):
        class_distributions.append(
            {
                'share_class_id': class_id,
                'share_class_name': class_row['share_class_name'],
                'shares': class_row['shares'],
                'preference_paid': _serialize_money(class_row['preference_paid']),
                'residual_paid': _serialize_money(class_row['residual_paid']),
                'total_paid': _serialize_money(class_totals[class_id]),
            }
        )

    for row in positions:
        class_total = class_totals.get(row['share_class_id'], Decimal('0'))
        class_share_count = Decimal(class_rows[row['share_class_id']]['shares'] or 1)
        payout = _money(class_total * Decimal(row['post_shares']) / class_share_count)
        holder_distributions.append(
            {
                'holder_id': row['holder_id'],
                'holder_name': row['holder_name'],
                'share_class_name': row['share_class_name'],
                'shares': int(row['post_shares']),
                'payout': _serialize_money(payout),
            }
        )

    holder_distributions.sort(key=lambda row: (_to_decimal(row['payout']), row['holder_name']), reverse=True)
    return {
        'exit_value': _serialize_money(exit_value),
        'remaining_exit': _serialize_money(residual_remaining),
        'class_distributions': class_distributions,
        'holder_distributions': holder_distributions,
    }


def simulate_financing_scenario(workspace, scenario_input: dict) -> dict:
    snapshot = build_cap_table_snapshot(workspace)
    positions, financing_meta = _build_post_financing_positions(snapshot, scenario_input)

    exit_values = scenario_input.get('exit_values') or []
    waterfall_results = [
        _waterfall_from_positions(positions, snapshot, _money(exit_value), financing_meta)
        for exit_value in exit_values
        if _to_decimal(exit_value) > 0
    ]

    positions.sort(key=lambda row: (-_to_decimal(row['ownership_percent']), row['holder_name'], row['share_class_name']))
    post_cap_table = [
        {
            'holder_id': row['holder_id'],
            'holder_name': row['holder_name'],
            'share_class_id': row['share_class_id'],
            'share_class_name': row['share_class_name'],
            'shares': int(row['post_shares']),
            'ownership_percent': _serialize_ratio(row['ownership_percent']),
            'new_money_shares': int(row['new_money_shares']),
            'pro_rata_shares': int(row['pro_rata_shares']),
            'anti_dilution_shares': int(row['anti_dilution_shares']),
            'option_pool_shares': int(row['option_pool_shares']),
            'scenario_source': row['scenario_source'],
        }
        for row in positions
    ]

    pre_cap_table = [
        {
            'holder_id': row['holder_id'],
            'holder_name': row['holder_name'],
            'share_class_id': row['share_class_id'],
            'share_class_name': row['share_class_name'],
            'shares': int(row['fully_diluted_shares']),
            'ownership_percent': _serialize_ratio(row['ownership_percent']),
            'scenario_source': row['scenario_source'],
        }
        for row in snapshot['holder_rows']
    ]

    return {
        'baseline': {
            'outstanding_shares': snapshot['outstanding_shares'],
            'fully_diluted_shares': snapshot['fully_diluted_shares'],
            'invested_capital': _serialize_money(snapshot['invested_capital']),
            'share_classes': [
                {
                    'id': row['id'],
                    'name': row['name'],
                    'class_type': row['class_type'],
                    'fully_diluted_shares': row['fully_diluted_shares'],
                    'ownership_percent': _serialize_ratio(row['ownership_percent']),
                    'invested_capital': _serialize_money(row['invested_capital']),
                }
                for row in snapshot['class_rows']
            ],
        },
        'financing': financing_meta,
        'pre_cap_table': pre_cap_table,
        'post_cap_table': post_cap_table,
        'waterfalls': waterfall_results,
    }


def create_scenario_report(workspace, scenario_input: dict, generated_by, *, title: str, reporting_period: str = ''):
    analysis = simulate_financing_scenario(workspace, scenario_input)
    report = EquityReport.objects.create(
        workspace=workspace,
        title=title,
        report_type='scenario_model',
        reporting_period=reporting_period,
        status='ready',
        payload={
            'scenario_input': scenario_input,
            'analysis': analysis,
        },
        generated_by=generated_by,
    )
    return report


def create_scenario_approval_request(workspace, scenario_input: dict, requested_by, *, title: str, reporting_period: str = ''):
    analysis = simulate_financing_scenario(workspace, scenario_input)
    policy = get_or_create_scenario_approval_policy(workspace)
    _validate_submission_policy(policy)
    approval = EquityScenarioApproval.objects.create(
        workspace=workspace,
        title=title,
        reporting_period=reporting_period,
        scenario_payload=scenario_input,
        analysis_payload=analysis,
        requested_by=requested_by,
        status=ScenarioApprovalStatus.PENDING,
        board_status=ScenarioReviewStatus.PENDING,
        legal_status=ScenarioReviewStatus.PENDING,
    )
    _apply_policy_sla_dates(approval, policy)
    approval.save(update_fields=['board_due_at', 'legal_due_at', 'updated_at'])
    _log_approval_event(
        approval,
        ScenarioApprovalEventType.SUBMITTED,
        'Scenario submitted for approval',
        'Scenario approval request was created and routed for board and legal review.',
        actor=requested_by,
        metadata={'reporting_period': reporting_period},
    )
    sync_equity_scenario_approval_task(approval)
    log_platform_audit_event(
        domain='equity',
        actor=requested_by,
        organization=workspace.organization,
        entity=workspace,
        event_type='equity_scenario_approval.created',
        action='approval_requested',
        resource_type='EquityScenarioApproval',
        resource_id=str(approval.id),
        subject_type='equity_scenario_approval',
        subject_id=str(approval.id),
        resource_name=approval.title,
        summary=f'Created equity approval request: {approval.title}',
        context={'board_due_at': approval.board_due_at.isoformat() if approval.board_due_at else None, 'legal_due_at': approval.legal_due_at.isoformat() if approval.legal_due_at else None},
        metadata={'reporting_period': reporting_period},
    )
    notify_scenario_requested(approval, get_board_review_users(workspace), get_legal_review_users(workspace))
    return approval


def _refresh_approval_status(approval: EquityScenarioApproval):
    if approval.board_status == ScenarioReviewStatus.REJECTED or approval.legal_status == ScenarioReviewStatus.REJECTED:
        approval.status = ScenarioApprovalStatus.REJECTED
    elif approval.board_status == ScenarioReviewStatus.APPROVED and approval.legal_status == ScenarioReviewStatus.APPROVED:
        approval.status = ScenarioApprovalStatus.APPROVED
    else:
        approval.status = ScenarioApprovalStatus.PENDING


@transaction.atomic
def approve_scenario_for_board(approval: EquityScenarioApproval, approver, comments: str = ''):
    if not user_can_board_approve(approval.workspace, approver):
        raise PermissionError('User is not authorized for board approval on this scenario.')
    approval.board_status = ScenarioReviewStatus.APPROVED
    approval.board_approved_by = approver
    approval.board_decided_at = timezone.now()
    if comments:
        approval.rejection_reason = comments
    _refresh_approval_status(approval)
    approval.save(update_fields=['board_status', 'board_approved_by', 'board_decided_at', 'status', 'rejection_reason', 'updated_at'])
    _log_approval_event(
        approval,
        ScenarioApprovalEventType.BOARD_APPROVED,
        'Board approval recorded',
        comments or 'Board reviewer approved the scenario.',
        actor=approver,
    )
    sync_equity_scenario_approval_task(approval)
    log_platform_audit_event(
        domain='equity',
        actor=approver,
        organization=approval.workspace.organization,
        entity=approval.workspace,
        event_type='equity_scenario_approval.board_approved',
        action='approval_progressed',
        resource_type='EquityScenarioApproval',
        resource_id=str(approval.id),
        subject_type='equity_scenario_approval',
        subject_id=str(approval.id),
        resource_name=approval.title,
        summary=f'Board approval recorded for {approval.title}',
        context={'status': approval.status, 'board_status': approval.board_status, 'legal_status': approval.legal_status},
        metadata={'comments': comments},
    )
    notify_scenario_review_decision(approval, reviewer_type='board', approved=True)
    return approval


@transaction.atomic
def approve_scenario_for_legal(approval: EquityScenarioApproval, approver, comments: str = ''):
    if not user_can_legal_approve(approval.workspace, approver):
        raise PermissionError('User is not authorized for legal approval on this scenario.')
    approval.legal_status = ScenarioReviewStatus.APPROVED
    approval.legal_approved_by = approver
    approval.legal_decided_at = timezone.now()
    if comments:
        approval.rejection_reason = comments
    _refresh_approval_status(approval)
    approval.save(update_fields=['legal_status', 'legal_approved_by', 'legal_decided_at', 'status', 'rejection_reason', 'updated_at'])
    _log_approval_event(
        approval,
        ScenarioApprovalEventType.LEGAL_APPROVED,
        'Legal approval recorded',
        comments or 'Legal reviewer approved the scenario.',
        actor=approver,
    )
    sync_equity_scenario_approval_task(approval)
    log_platform_audit_event(
        domain='equity',
        actor=approver,
        organization=approval.workspace.organization,
        entity=approval.workspace,
        event_type='equity_scenario_approval.legal_approved',
        action='approval_progressed',
        resource_type='EquityScenarioApproval',
        resource_id=str(approval.id),
        subject_type='equity_scenario_approval',
        subject_id=str(approval.id),
        resource_name=approval.title,
        summary=f'Legal approval recorded for {approval.title}',
        context={'status': approval.status, 'board_status': approval.board_status, 'legal_status': approval.legal_status},
        metadata={'comments': comments},
    )
    notify_scenario_review_decision(approval, reviewer_type='legal', approved=True)
    return approval


@transaction.atomic
def reject_scenario_approval(approval: EquityScenarioApproval, reviewer_type: str, reviewer, comments: str = ''):
    if reviewer_type == 'board' and not user_can_board_approve(approval.workspace, reviewer):
        raise PermissionError('User is not authorized for board rejection on this scenario.')
    if reviewer_type == 'legal' and not user_can_legal_approve(approval.workspace, reviewer):
        raise PermissionError('User is not authorized for legal rejection on this scenario.')
    if reviewer_type == 'board':
        approval.board_status = ScenarioReviewStatus.REJECTED
        approval.board_approved_by = reviewer
        approval.board_decided_at = timezone.now()
        update_fields = ['board_status', 'board_approved_by', 'board_decided_at']
    else:
        approval.legal_status = ScenarioReviewStatus.REJECTED
        approval.legal_approved_by = reviewer
        approval.legal_decided_at = timezone.now()
        update_fields = ['legal_status', 'legal_approved_by', 'legal_decided_at']
    approval.rejection_reason = comments or approval.rejection_reason
    _refresh_approval_status(approval)
    approval.save(update_fields=[*update_fields, 'rejection_reason', 'status', 'updated_at'])
    _log_approval_event(
        approval,
        ScenarioApprovalEventType.REJECTED,
        f'{reviewer_type.title()} rejection recorded',
        comments or f'{reviewer_type.title()} reviewer rejected the scenario.',
        actor=reviewer,
        metadata={'reviewer_type': reviewer_type},
    )
    sync_equity_scenario_approval_task(approval)
    log_platform_audit_event(
        domain='equity',
        actor=reviewer,
        organization=approval.workspace.organization,
        entity=approval.workspace,
        event_type='equity_scenario_approval.rejected',
        action='approval_rejected',
        resource_type='EquityScenarioApproval',
        resource_id=str(approval.id),
        subject_type='equity_scenario_approval',
        subject_id=str(approval.id),
        resource_name=approval.title,
        summary=f'{reviewer_type.title()} rejected equity approval: {approval.title}',
        context={'status': approval.status, 'reviewer_type': reviewer_type},
        metadata={'comments': comments},
    )
    notify_scenario_review_decision(approval, reviewer_type=reviewer_type, approved=False)
    return approval


def _get_or_create_scenario_security(workspace, scenario_input: dict):
    share_class_id = scenario_input.get('share_class')
    if share_class_id:
        return EquityShareClass.objects.get(pk=share_class_id, workspace=workspace)
    name = scenario_input.get('name') or f'Preferred Round {timezone.now().date().isoformat()}'
    share_class, _ = EquityShareClass.objects.get_or_create(
        workspace=workspace,
        name=name,
        defaults={
            'class_type': ShareClassType.PREFERRED,
            'authorized_shares': 0,
            'issued_shares': 0,
            'liquidation_preference': '1x non-participating',
            'preference_multiple': Decimal('1'),
            'participating_preference': False,
            'participation_cap_multiple': None,
            'liquidation_seniority': 1,
            'conversion_price': _to_decimal(scenario_input.get('price_per_share') or 0),
            'pro_rata_rights': True,
        },
    )
    return share_class


def _get_or_create_shareholder(workspace, holder_name: str, shareholder_type: str = ShareholderType.INVESTOR):
    shareholder, _ = EquityShareholder.objects.get_or_create(
        workspace=workspace,
        name=holder_name,
        email='',
        defaults={'shareholder_type': shareholder_type},
    )
    if not shareholder.shareholder_type:
        shareholder.shareholder_type = shareholder_type
        shareholder.save(update_fields=['shareholder_type', 'updated_at'])
    return shareholder


def _recalculate_holding_percentages(workspace):
    holdings = list(EquityHolding.objects.filter(workspace=workspace).select_related('share_class'))
    total_diluted = sum(int(item.diluted_quantity or item.quantity or 0) for item in holdings) or 1
    issued_by_class = {}
    for holding in holdings:
        holding.ownership_percent = (Decimal(int(holding.diluted_quantity or holding.quantity or 0)) * Decimal('100') / Decimal(total_diluted)).quantize(Decimal('0.0001'))
        holding.save(update_fields=['ownership_percent', 'updated_at'])
        issued_by_class[holding.share_class_id] = issued_by_class.get(holding.share_class_id, 0) + int(holding.quantity or 0)
    for share_class in EquityShareClass.objects.filter(workspace=workspace):
        issued_shares = issued_by_class.get(share_class.id, 0)
        updates = []
        if share_class.issued_shares != issued_shares:
            share_class.issued_shares = issued_shares
            updates.append('issued_shares')
        if share_class.authorized_shares < share_class.issued_shares:
            share_class.authorized_shares = share_class.issued_shares
            updates.append('authorized_shares')
        if updates:
            share_class.save(update_fields=[*updates, 'updated_at'])


@transaction.atomic
def commit_financing_scenario(workspace, scenario_input: dict, actor):
    analysis = simulate_financing_scenario(workspace, scenario_input)
    financing = analysis['financing']
    price_per_share = _to_decimal(financing['price_per_share'])
    scenario_security = _get_or_create_scenario_security(workspace, scenario_input)

    funding_round = EquityFundingRound.objects.create(
        workspace=workspace,
        name=scenario_input.get('name') or scenario_security.name,
        instrument_type='equity',
        share_class=scenario_security,
        announced_at=timezone.now().date(),
        pre_money_valuation=_to_decimal(financing['pre_money_valuation']),
        post_money_valuation=_to_decimal(financing['post_money_valuation']),
        amount_raised=_to_decimal(scenario_input.get('amount_raised') or 0),
        price_per_share=price_per_share,
        new_shares_issued=int(financing['new_money_shares']),
        option_pool_top_up=int(financing['option_pool_top_up']),
        apply_pro_rata=bool(scenario_input.get('apply_pro_rata', True)),
        scenario_assumptions={'scenario_input': scenario_input, 'analysis': analysis},
        notes='Committed from scenario modeling workflow.',
    )

    committed_holders = []
    post_rows = analysis['post_cap_table']
    for row in post_rows:
        if row['scenario_source'] not in {'holding', 'pro_rata', 'new_investor'} and int(row.get('anti_dilution_shares') or 0) <= 0:
            continue

        if row['holder_id'] == 'scenario-new-investor':
            shareholder = _get_or_create_shareholder(workspace, row['holder_name'], ShareholderType.INVESTOR)
        elif row['holder_id'].startswith('grant-pool-'):
            continue
        else:
            shareholder = EquityShareholder.objects.get(pk=row['holder_id'], workspace=workspace)

        share_class = scenario_security if row['share_class_id'] == financing['security']['id'] else EquityShareClass.objects.get(pk=row['share_class_id'], workspace=workspace)
        holding, _ = EquityHolding.objects.get_or_create(
            workspace=workspace,
            shareholder=shareholder,
            share_class=share_class,
            defaults={
                'quantity': 0,
                'diluted_quantity': 0,
                'issue_price_per_share': price_per_share if share_class == scenario_security else Decimal('0'),
                'invested_amount': Decimal('0'),
                'pro_rata_eligible': bool(share_class.pro_rata_rights),
                'pro_rata_take_up_percent': Decimal('100'),
                'issued_at': timezone.now().date(),
            },
        )

        issued_delta = int(row.get('new_money_shares') or 0)
        anti_dilution_delta = int(row.get('anti_dilution_shares') or 0)
        if issued_delta:
            holding.quantity += issued_delta
            holding.diluted_quantity += issued_delta
            holding.issue_price_per_share = price_per_share
            holding.invested_amount = _money(holding.invested_amount + (price_per_share * Decimal(issued_delta)))
        if anti_dilution_delta:
            holding.diluted_quantity += anti_dilution_delta
        holding.pro_rata_eligible = bool(holding.pro_rata_eligible or share_class.pro_rata_rights)
        holding.save()

        if issued_delta:
            EquityTransaction.objects.create(
                workspace=workspace,
                transaction_type=TransactionType.ISSUE,
                shareholder=shareholder,
                share_class=share_class,
                quantity=issued_delta,
                price_per_share=price_per_share,
                effective_date=timezone.now().date(),
                approval_status=ApprovalStatus.EXECUTED,
                compliance_checked=True,
                digital_signature_required=False,
                audit_metadata={'source': 'scenario_commit', 'funding_round': str(funding_round.id)},
                created_by=actor,
            )
        if anti_dilution_delta:
            EquityTransaction.objects.create(
                workspace=workspace,
                transaction_type=TransactionType.CONVERSION,
                shareholder=shareholder,
                share_class=share_class,
                quantity=anti_dilution_delta,
                price_per_share=Decimal('0'),
                effective_date=timezone.now().date(),
                approval_status=ApprovalStatus.EXECUTED,
                compliance_checked=True,
                digital_signature_required=False,
                audit_metadata={'source': 'scenario_commit_anti_dilution', 'funding_round': str(funding_round.id)},
                created_by=actor,
            )
        committed_holders.append({'holder_name': shareholder.name, 'share_class_name': share_class.name, 'issued_delta': issued_delta, 'anti_dilution_delta': anti_dilution_delta})

    option_pool_top_up = int(financing.get('option_pool_top_up') or 0)
    if option_pool_top_up > 0:
        option_class, _ = EquityShareClass.objects.get_or_create(
            workspace=workspace,
            name='Option Pool',
            defaults={
                'class_type': ShareClassType.ESOP,
                'authorized_shares': option_pool_top_up,
                'issued_shares': 0,
            },
        )
        if option_class.authorized_shares < option_pool_top_up:
            option_class.authorized_shares += option_pool_top_up
            option_class.save(update_fields=['authorized_shares', 'updated_at'])
        reserve_record, _ = EquityOptionPoolReserve.objects.get_or_create(
            workspace=workspace,
            share_class=option_class,
            funding_round=funding_round,
            defaults={'reserved_shares': 0, 'allocated_shares': 0, 'notes': 'Scenario committed option pool reserve.'},
        )
        reserve_record.reserved_shares += option_pool_top_up
        reserve_record.funding_round = funding_round
        reserve_record.save(update_fields=['reserved_shares', 'funding_round', 'updated_at'])

    _recalculate_holding_percentages(workspace)
    return {
        'funding_round_id': str(funding_round.id),
        'funding_round_name': funding_round.name,
        'security': scenario_security.name,
        'committed_holders': committed_holders,
        'analysis': analysis,
    }


@transaction.atomic
def commit_approved_scenario(approval: EquityScenarioApproval, actor):
    if approval.status != ScenarioApprovalStatus.APPROVED:
        raise ValueError('Scenario approval is not ready to commit.')
    result = commit_financing_scenario(approval.workspace, approval.scenario_payload, actor)
    funding_round = EquityFundingRound.objects.get(pk=result['funding_round_id'])
    approval.status = ScenarioApprovalStatus.COMMITTED
    approval.committed_round = funding_round
    approval.save(update_fields=['status', 'committed_round', 'updated_at'])
    _log_approval_event(
        approval,
        ScenarioApprovalEventType.COMMITTED,
        'Scenario committed to cap table',
        f'Scenario was committed as {funding_round.name}.',
        actor=actor,
        metadata={'funding_round_id': str(funding_round.id)},
    )
    notify_scenario_committed(approval)
    result['approval_id'] = str(approval.id)
    return result


def get_scenario_overview(workspace) -> dict:
    snapshot = build_cap_table_snapshot(workspace)
    latest_round = workspace.equity_funding_rounds.order_by('-announced_at', '-created_at').select_related('share_class').first()
    latest_valuation = workspace.equity_valuations.order_by('-valuation_date', '-created_at').first()
    return {
        'cap_table': {
            'outstanding_shares': snapshot['outstanding_shares'],
            'fully_diluted_shares': snapshot['fully_diluted_shares'],
            'invested_capital': _serialize_money(snapshot['invested_capital']),
            'holders': [
                {
                    'holder_name': row['holder_name'],
                    'share_class_name': row['share_class_name'],
                    'shares': int(row['fully_diluted_shares']),
                    'ownership_percent': _serialize_ratio(row['ownership_percent']),
                }
                for row in snapshot['holder_rows']
            ],
        },
        'defaults': {
            'latest_round_name': latest_round.name if latest_round else '',
            'latest_round_pre_money': _serialize_money(latest_round.pre_money_valuation if latest_round else 0),
            'latest_round_price_per_share': _serialize_ratio(latest_round.price_per_share if latest_round else 0),
            'latest_valuation_price_per_share': _serialize_ratio(latest_valuation.price_per_share if latest_valuation else 0),
            'default_exit_values': ['25000000.00', '50000000.00', '100000000.00'],
        },
    }