import uuid

from django.conf import settings
from django.db import models

from atonixcorp.models import ChartOfAccounts, Entity, EntityStaff, JournalEntry, TaxCalculation


class EquityWorkspaceType(models.TextChoices):
    STANDALONE = 'standalone', 'Standalone'
    ACCOUNTING = 'accounting', 'Accounting'
    EQUITY = 'equity', 'Equity'
    COMBINED = 'combined', 'Combined'


class ShareholderType(models.TextChoices):
    INDIVIDUAL = 'individual', 'Individual'
    ENTITY = 'entity', 'Entity'
    EMPLOYEE = 'employee', 'Employee'
    INVESTOR = 'investor', 'Investor'


class ShareClassType(models.TextChoices):
    COMMON = 'common', 'Common'
    PREFERRED = 'preferred', 'Preferred'
    ESOP = 'esop', 'ESOP'
    WARRANT = 'warrant', 'Warrant'
    SAFE = 'safe', 'SAFE'
    CONVERTIBLE = 'convertible', 'Convertible'


class AntiDilutionType(models.TextChoices):
    NONE = 'none', 'None'
    FULL_RATCHET = 'full_ratchet', 'Full Ratchet'
    WEIGHTED_AVERAGE = 'weighted_average', 'Weighted Average'


class AntiDilutionBasis(models.TextChoices):
    BROAD_BASED = 'broad_based', 'Broad-Based'
    NARROW_BASED = 'narrow_based', 'Narrow-Based'


class InstrumentType(models.TextChoices):
    EQUITY = 'equity', 'Equity'
    SAFE = 'safe', 'SAFE'
    CONVERTIBLE_NOTE = 'convertible_note', 'Convertible Note'
    WARRANT = 'warrant', 'Warrant'
    OPTION = 'option', 'Option'


class TransactionType(models.TextChoices):
    ISSUE = 'issue', 'Issue'
    TRANSFER = 'transfer', 'Transfer'
    EXERCISE = 'exercise', 'Exercise'
    CANCELLATION = 'cancellation', 'Cancellation'
    CONVERSION = 'conversion', 'Conversion'


class ApprovalStatus(models.TextChoices):
    DRAFT = 'draft', 'Draft'
    PENDING = 'pending', 'Pending'
    APPROVED = 'approved', 'Approved'
    REJECTED = 'rejected', 'Rejected'
    EXECUTED = 'executed', 'Executed'


class ValuationMethod(models.TextChoices):
    DCF = 'dcf', 'DCF'
    COMPS = 'comps', 'Comparables'
    MARKET = 'market', 'Market'
    BOARD = 'board', 'Board Approved'


class EquityGrantType(models.TextChoices):
    FOUNDER = 'founder', 'Founder Vesting'
    STOCK_OPTION = 'stock_option', 'Stock Option'
    ESOP = 'esop', 'ESOP'
    RSU = 'rsu', 'RSU'
    RESTRICTED_STOCK = 'restricted_stock', 'Restricted Stock'
    WARRANT = 'warrant', 'Warrant'


class VestingInterval(models.TextChoices):
    MONTHLY = 'monthly', 'Monthly'
    QUARTERLY = 'quarterly', 'Quarterly'
    ANNUAL = 'annual', 'Annual'
    CUSTOM = 'custom', 'Custom'


class GrantLifecycleStatus(models.TextChoices):
    GRANTED = 'granted', 'Granted'
    ACTIVE = 'active', 'Active'
    FULLY_VESTED = 'fully_vested', 'Fully Vested'
    EXERCISED = 'exercised', 'Exercised'
    EXPIRED = 'expired', 'Expired'
    FORFEITED = 'forfeited', 'Forfeited'
    CANCELLED = 'cancelled', 'Cancelled'


class AccelerationType(models.TextChoices):
    NONE = 'none', 'None'
    SINGLE = 'single', 'Single Trigger'
    DOUBLE = 'double', 'Double Trigger'


class TerminationTreatment(models.TextChoices):
    FORFEIT_UNVESTED = 'forfeit_unvested', 'Forfeit Unvested'
    CONTINUE_VESTING = 'continue_vesting', 'Continue Standard Vesting'
    ACCELERATE_TO_CLIFF = 'accelerate_to_cliff', 'Accelerate To Cliff'
    FULL_ACCELERATION = 'full_acceleration', 'Full Acceleration'


class VestingEventType(models.TextChoices):
    SCHEDULED = 'scheduled', 'Scheduled Vesting'
    CLIFF_RELEASE = 'cliff_release', 'Cliff Release'
    ACCELERATION = 'acceleration', 'Acceleration'
    FORFEITURE = 'forfeiture', 'Forfeiture'
    EXERCISE = 'exercise', 'Exercise'
    MANUAL = 'manual', 'Manual Adjustment'


class VestingEventStatus(models.TextChoices):
    PENDING = 'pending', 'Pending'
    VESTED = 'vested', 'Vested'
    FORFEITED = 'forfeited', 'Forfeited'
    EXERCISED = 'exercised', 'Exercised'
    CANCELLED = 'cancelled', 'Cancelled'


class ExercisePaymentMethod(models.TextChoices):
    BANK_TRANSFER = 'bank_transfer', 'Bank Transfer'
    PAYROLL_DEDUCTION = 'payroll_deduction', 'Payroll Deduction'
    WALLET = 'wallet', 'Wallet'
    CASHLESS = 'cashless', 'Cashless Exercise'


class ExerciseRequestStatus(models.TextChoices):
    REQUESTED = 'requested', 'Requested'
    FINANCE_REVIEW = 'finance_review', 'Finance Review'
    LEGAL_REVIEW = 'legal_review', 'Legal Review'
    APPROVED = 'approved', 'Approved'
    REJECTED = 'rejected', 'Rejected'
    COMPLETED = 'completed', 'Completed'
    CANCELLED = 'cancelled', 'Cancelled'


class ApprovalDecisionStatus(models.TextChoices):
    PENDING = 'pending', 'Pending'
    APPROVED = 'approved', 'Approved'
    REJECTED = 'rejected', 'Rejected'
    SKIPPED = 'skipped', 'Skipped'


class PaymentStatus(models.TextChoices):
    PENDING = 'pending', 'Pending'
    PROCESSING = 'processing', 'Processing'
    PAID = 'paid', 'Paid'
    FAILED = 'failed', 'Failed'
    WAIVED = 'waived', 'Waived'


class CertificateStatus(models.TextChoices):
    DRAFT = 'draft', 'Draft'
    ISSUED = 'issued', 'Issued'
    VOID = 'void', 'Void'


class PayrollSyncStatus(models.TextChoices):
    QUEUED = 'queued', 'Queued'
    SYNCED = 'synced', 'Synced'
    FAILED = 'failed', 'Failed'


class ExternalAdapterType(models.TextChoices):
    PAYROLL = 'payroll', 'Payroll'
    PAYMENT = 'payment', 'Payment'


class DeliveryChannel(models.TextChoices):
    IN_APP = 'in_app', 'In-App'
    EMAIL = 'email', 'Email'
    DOCUMENT = 'document', 'Document'
    WEBHOOK = 'webhook', 'Webhook'


class DeliveryStatus(models.TextChoices):
    QUEUED = 'queued', 'Queued'
    SENT = 'sent', 'Sent'
    FAILED = 'failed', 'Failed'


class ScenarioReviewStatus(models.TextChoices):
    PENDING = 'pending', 'Pending'
    APPROVED = 'approved', 'Approved'
    REJECTED = 'rejected', 'Rejected'


class ScenarioApprovalStatus(models.TextChoices):
    PENDING = 'pending', 'Pending'
    APPROVED = 'approved', 'Approved'
    REJECTED = 'rejected', 'Rejected'
    COMMITTED = 'committed', 'Committed'


class ScenarioApprovalEventType(models.TextChoices):
    SUBMITTED = 'submitted', 'Submitted'
    BOARD_APPROVED = 'board_approved', 'Board Approved'
    LEGAL_APPROVED = 'legal_approved', 'Legal Approved'
    REJECTED = 'rejected', 'Rejected'
    REMINDER = 'reminder', 'Reminder Sent'
    ESCALATED = 'escalated', 'Escalated'
    COMMITTED = 'committed', 'Committed'


class WorkspaceEquityProfile(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.OneToOneField(Entity, on_delete=models.CASCADE, related_name='equity_profile')
    workspace_type = models.CharField(max_length=20, choices=EquityWorkspaceType.choices, default=EquityWorkspaceType.ACCOUNTING)
    equity_enabled = models.BooleanField(default=False)
    ownership_registry_enabled = models.BooleanField(default=False)
    cap_table_enabled = models.BooleanField(default=False)
    valuation_enabled = models.BooleanField(default=False)
    equity_transactions_enabled = models.BooleanField(default=False)
    governance_reporting_enabled = models.BooleanField(default=False)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return f'{self.workspace.name} equity profile'


class EquityScenarioApprovalPolicy(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.OneToOneField(Entity, on_delete=models.CASCADE, related_name='equity_scenario_approval_policy')
    board_reviewers = models.ManyToManyField(EntityStaff, blank=True, related_name='equity_board_review_policies')
    legal_reviewers = models.ManyToManyField(EntityStaff, blank=True, related_name='equity_legal_review_policies')
    board_escalation_reviewers = models.ManyToManyField(EntityStaff, blank=True, related_name='equity_board_escalation_policies')
    legal_escalation_reviewers = models.ManyToManyField(EntityStaff, blank=True, related_name='equity_legal_escalation_policies')
    require_explicit_reviewers = models.BooleanField(default=False)
    require_designated_backups = models.BooleanField(default=False)
    board_sla_hours = models.PositiveIntegerField(default=72)
    legal_sla_hours = models.PositiveIntegerField(default=72)
    escalation_enabled = models.BooleanField(default=True)
    escalation_grace_hours = models.PositiveIntegerField(default=24)
    reminder_frequency_hours = models.PositiveIntegerField(default=24)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return f'{self.workspace.name} scenario approval policy'


class EquityShareholder(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_shareholders')
    name = models.CharField(max_length=255)
    shareholder_type = models.CharField(max_length=20, choices=ShareholderType.choices, default=ShareholderType.INDIVIDUAL)
    email = models.EmailField(blank=True, default='')
    beneficial_owner = models.BooleanField(default=False)
    voting_rights_percent = models.DecimalField(max_digits=8, decimal_places=2, default=0)
    kyc_status = models.CharField(max_length=50, blank=True, default='pending')
    aml_status = models.CharField(max_length=50, blank=True, default='pending')
    notes = models.TextField(blank=True, default='')
    created_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='created_equity_shareholders')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['name']
        unique_together = ('workspace', 'name', 'email')

    def __str__(self):
        return self.name


class EquityShareClass(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_share_classes')
    name = models.CharField(max_length=255)
    class_type = models.CharField(max_length=20, choices=ShareClassType.choices, default=ShareClassType.COMMON)
    authorized_shares = models.BigIntegerField(default=0)
    issued_shares = models.BigIntegerField(default=0)
    liquidation_preference = models.CharField(max_length=255, blank=True, default='')
    preference_multiple = models.DecimalField(max_digits=8, decimal_places=4, default=0)
    participating_preference = models.BooleanField(default=False)
    participation_cap_multiple = models.DecimalField(max_digits=8, decimal_places=4, null=True, blank=True)
    liquidation_seniority = models.PositiveIntegerField(default=0)
    conversion_price = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    anti_dilution_type = models.CharField(max_length=30, choices=AntiDilutionType.choices, default=AntiDilutionType.NONE)
    anti_dilution_basis = models.CharField(max_length=30, choices=AntiDilutionBasis.choices, default=AntiDilutionBasis.BROAD_BASED)
    pro_rata_rights = models.BooleanField(default=False)
    voting_rights = models.BooleanField(default=True)
    par_value = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    currency = models.CharField(max_length=10, default='USD')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['name']
        unique_together = ('workspace', 'name')

    def __str__(self):
        return f'{self.workspace.name} - {self.name}'


class EquityHolding(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_holdings')
    shareholder = models.ForeignKey(EquityShareholder, on_delete=models.CASCADE, related_name='holdings')
    share_class = models.ForeignKey(EquityShareClass, on_delete=models.CASCADE, related_name='holdings')
    quantity = models.BigIntegerField(default=0)
    diluted_quantity = models.BigIntegerField(default=0)
    ownership_percent = models.DecimalField(max_digits=8, decimal_places=4, default=0)
    issued_at = models.DateField(null=True, blank=True)
    vesting_start = models.DateField(null=True, blank=True)
    vesting_end = models.DateField(null=True, blank=True)
    issue_price_per_share = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    invested_amount = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    pro_rata_eligible = models.BooleanField(default=False)
    pro_rata_take_up_percent = models.DecimalField(max_digits=8, decimal_places=2, default=100)
    strike_price = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-updated_at']
        unique_together = ('workspace', 'shareholder', 'share_class')


class EquityFundingRound(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_funding_rounds')
    name = models.CharField(max_length=255)
    instrument_type = models.CharField(max_length=30, choices=InstrumentType.choices, default=InstrumentType.EQUITY)
    share_class = models.ForeignKey(EquityShareClass, null=True, blank=True, on_delete=models.SET_NULL, related_name='funding_rounds')
    announced_at = models.DateField(null=True, blank=True)
    pre_money_valuation = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    post_money_valuation = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    amount_raised = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    price_per_share = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    new_shares_issued = models.BigIntegerField(default=0)
    option_pool_top_up = models.BigIntegerField(default=0)
    apply_pro_rata = models.BooleanField(default=True)
    scenario_assumptions = models.JSONField(default=dict, blank=True)
    notes = models.TextField(blank=True, default='')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-announced_at', '-created_at']


class EquityOptionPoolReserve(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_option_pool_reserves')
    share_class = models.ForeignKey(EquityShareClass, on_delete=models.CASCADE, related_name='option_pool_reserves')
    funding_round = models.ForeignKey(EquityFundingRound, null=True, blank=True, on_delete=models.SET_NULL, related_name='option_pool_reserves')
    reserved_shares = models.BigIntegerField(default=0)
    allocated_shares = models.BigIntegerField(default=0)
    notes = models.TextField(blank=True, default='')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-updated_at']


class EquityScenarioApproval(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_scenario_approvals')
    title = models.CharField(max_length=255)
    reporting_period = models.CharField(max_length=100, blank=True, default='')
    scenario_payload = models.JSONField(default=dict, blank=True)
    analysis_payload = models.JSONField(default=dict, blank=True)
    status = models.CharField(max_length=20, choices=ScenarioApprovalStatus.choices, default=ScenarioApprovalStatus.PENDING)
    board_status = models.CharField(max_length=20, choices=ScenarioReviewStatus.choices, default=ScenarioReviewStatus.PENDING)
    legal_status = models.CharField(max_length=20, choices=ScenarioReviewStatus.choices, default=ScenarioReviewStatus.PENDING)
    board_due_at = models.DateTimeField(null=True, blank=True)
    legal_due_at = models.DateTimeField(null=True, blank=True)
    board_last_reminder_at = models.DateTimeField(null=True, blank=True)
    legal_last_reminder_at = models.DateTimeField(null=True, blank=True)
    board_escalated_at = models.DateTimeField(null=True, blank=True)
    legal_escalated_at = models.DateTimeField(null=True, blank=True)
    requested_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='requested_equity_scenario_approvals')
    board_approved_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='board_approved_equity_scenarios')
    legal_approved_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='legal_approved_equity_scenarios')
    board_decided_at = models.DateTimeField(null=True, blank=True)
    legal_decided_at = models.DateTimeField(null=True, blank=True)
    rejection_reason = models.TextField(blank=True, default='')
    committed_round = models.ForeignKey(EquityFundingRound, null=True, blank=True, on_delete=models.SET_NULL, related_name='scenario_approval_records')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-created_at']


class EquityScenarioApprovalEvent(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    approval = models.ForeignKey(EquityScenarioApproval, on_delete=models.CASCADE, related_name='events')
    actor = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='equity_scenario_approval_events')
    event_type = models.CharField(max_length=40, choices=ScenarioApprovalEventType.choices)
    title = models.CharField(max_length=255)
    message = models.TextField(blank=True, default='')
    metadata = models.JSONField(default=dict, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ['-created_at']


class EquityValuation(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_valuations')
    title = models.CharField(max_length=255)
    method = models.CharField(max_length=20, choices=ValuationMethod.choices, default=ValuationMethod.BOARD)
    valuation_date = models.DateField()
    enterprise_value = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    equity_value = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    price_per_share = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    market_notes = models.TextField(blank=True, default='')
    benchmark_data = models.JSONField(default=dict, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-valuation_date', '-created_at']


class EquityTransaction(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_transactions')
    transaction_type = models.CharField(max_length=20, choices=TransactionType.choices)
    shareholder = models.ForeignKey(EquityShareholder, null=True, blank=True, on_delete=models.SET_NULL, related_name='transactions')
    share_class = models.ForeignKey(EquityShareClass, null=True, blank=True, on_delete=models.SET_NULL, related_name='transactions')
    quantity = models.BigIntegerField(default=0)
    price_per_share = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    effective_date = models.DateField()
    approval_status = models.CharField(max_length=20, choices=ApprovalStatus.choices, default=ApprovalStatus.DRAFT)
    compliance_checked = models.BooleanField(default=False)
    digital_signature_required = models.BooleanField(default=True)
    audit_metadata = models.JSONField(default=dict, blank=True)
    created_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='created_equity_transactions')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-effective_date', '-created_at']


class EquityReport(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_reports')
    title = models.CharField(max_length=255)
    report_type = models.CharField(max_length=100)
    reporting_period = models.CharField(max_length=100, blank=True, default='')
    status = models.CharField(max_length=30, default='ready')
    payload = models.JSONField(default=dict, blank=True)
    generated_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='generated_equity_reports')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-created_at']


class EquityGrant(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_grants')
    grant_number = models.CharField(max_length=100)
    shareholder = models.ForeignKey(EquityShareholder, on_delete=models.CASCADE, related_name='grants')
    employee = models.ForeignKey(EntityStaff, null=True, blank=True, on_delete=models.SET_NULL, related_name='equity_grants')
    share_class = models.ForeignKey(EquityShareClass, on_delete=models.CASCADE, related_name='grants')
    grant_type = models.CharField(max_length=30, choices=EquityGrantType.choices, default=EquityGrantType.STOCK_OPTION)
    lifecycle_status = models.CharField(max_length=30, choices=GrantLifecycleStatus.choices, default=GrantLifecycleStatus.GRANTED)
    total_units = models.BigIntegerField(default=0)
    exercised_units = models.BigIntegerField(default=0)
    forfeited_units = models.BigIntegerField(default=0)
    exercise_price = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    grant_date = models.DateField()
    vesting_start_date = models.DateField()
    cliff_months = models.PositiveIntegerField(default=12)
    vesting_months = models.PositiveIntegerField(default=48)
    vesting_interval = models.CharField(max_length=20, choices=VestingInterval.choices, default=VestingInterval.MONTHLY)
    custom_schedule = models.JSONField(default=list, blank=True)
    acceleration_type = models.CharField(max_length=20, choices=AccelerationType.choices, default=AccelerationType.NONE)
    single_trigger_acceleration_percent = models.DecimalField(max_digits=5, decimal_places=2, default=100)
    double_trigger_acceleration_percent = models.DecimalField(max_digits=5, decimal_places=2, default=100)
    single_trigger_event_date = models.DateField(null=True, blank=True)
    double_trigger_event_date = models.DateField(null=True, blank=True)
    termination_treatment = models.CharField(max_length=30, choices=TerminationTreatment.choices, default=TerminationTreatment.FORFEIT_UNVESTED)
    post_termination_exercise_days = models.PositiveIntegerField(default=90)
    termination_date = models.DateField(null=True, blank=True)
    expiration_date = models.DateField(null=True, blank=True)
    fully_vested_at = models.DateField(null=True, blank=True)
    last_vesting_calculated_at = models.DateTimeField(null=True, blank=True)
    grant_package_file = models.FileField(upload_to='equity/grant_packages/', null=True, blank=True)
    notes = models.TextField(blank=True, default='')
    metadata = models.JSONField(default=dict, blank=True)
    created_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='created_equity_grants')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-grant_date', '-created_at']
        unique_together = ('workspace', 'grant_number')

    def __str__(self):
        return f'{self.grant_number} - {self.shareholder.name}'


class EquityVestingEvent(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_vesting_events')
    grant = models.ForeignKey(EquityGrant, on_delete=models.CASCADE, related_name='vesting_events')
    event_type = models.CharField(max_length=30, choices=VestingEventType.choices, default=VestingEventType.SCHEDULED)
    status = models.CharField(max_length=20, choices=VestingEventStatus.choices, default=VestingEventStatus.PENDING)
    vest_date = models.DateField()
    units = models.BigIntegerField(default=0)
    source_reference = models.CharField(max_length=120, blank=True, default='')
    trigger_name = models.CharField(max_length=120, blank=True, default='')
    notes = models.TextField(blank=True, default='')
    metadata = models.JSONField(default=dict, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['vest_date', 'created_at']


class EquityExerciseRequest(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_exercise_requests')
    grant = models.ForeignKey(EquityGrant, on_delete=models.CASCADE, related_name='exercise_requests')
    shareholder = models.ForeignKey(EquityShareholder, on_delete=models.CASCADE, related_name='exercise_requests')
    requested_units = models.BigIntegerField(default=0)
    approved_units = models.BigIntegerField(default=0)
    strike_price_per_unit = models.DecimalField(max_digits=14, decimal_places=4, default=0)
    strike_payment_amount = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    tax_withholding_amount = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    payment_method = models.CharField(max_length=30, choices=ExercisePaymentMethod.choices, default=ExercisePaymentMethod.BANK_TRANSFER)
    payment_status = models.CharField(max_length=20, choices=PaymentStatus.choices, default=PaymentStatus.PENDING)
    status = models.CharField(max_length=30, choices=ExerciseRequestStatus.choices, default=ExerciseRequestStatus.REQUESTED)
    requested_at = models.DateTimeField(auto_now_add=True)
    due_date = models.DateField(null=True, blank=True)
    completed_at = models.DateTimeField(null=True, blank=True)
    exercise_date = models.DateField(null=True, blank=True)
    notes = models.TextField(blank=True, default='')
    tax_calculation = models.ForeignKey(TaxCalculation, null=True, blank=True, on_delete=models.SET_NULL, related_name='equity_exercise_requests')
    journal_entry = models.ForeignKey(JournalEntry, null=True, blank=True, on_delete=models.SET_NULL, related_name='equity_exercise_requests')
    created_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='created_equity_exercise_requests')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-requested_at']


class EquityExerciseApproval(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    exercise_request = models.ForeignKey(EquityExerciseRequest, on_delete=models.CASCADE, related_name='approvals')
    approver = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, related_name='equity_exercise_approvals')
    approval_order = models.PositiveIntegerField(default=1)
    status = models.CharField(max_length=20, choices=ApprovalDecisionStatus.choices, default=ApprovalDecisionStatus.PENDING)
    decided_at = models.DateTimeField(null=True, blank=True)
    comments = models.TextField(blank=True, default='')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['approval_order', 'created_at']
        unique_together = ('exercise_request', 'approver', 'approval_order')


class EquityShareCertificate(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_certificates')
    exercise_request = models.OneToOneField(EquityExerciseRequest, null=True, blank=True, on_delete=models.SET_NULL, related_name='certificate')
    grant = models.ForeignKey(EquityGrant, on_delete=models.CASCADE, related_name='certificates')
    certificate_number = models.CharField(max_length=120)
    issued_to = models.ForeignKey(EquityShareholder, on_delete=models.CASCADE, related_name='certificates')
    share_class = models.ForeignKey(EquityShareClass, on_delete=models.CASCADE, related_name='certificates')
    issued_units = models.BigIntegerField(default=0)
    issue_date = models.DateField()
    status = models.CharField(max_length=20, choices=CertificateStatus.choices, default=CertificateStatus.DRAFT)
    certificate_payload = models.JSONField(default=dict, blank=True)
    pdf_file = models.FileField(upload_to='equity/certificates/', null=True, blank=True)
    issued_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='issued_equity_certificates')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-issue_date', '-created_at']
        unique_together = ('workspace', 'certificate_number')


class EquityPayrollTaxEvent(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_payroll_tax_events')
    grant = models.ForeignKey(EquityGrant, null=True, blank=True, on_delete=models.SET_NULL, related_name='payroll_tax_events')
    exercise_request = models.ForeignKey(EquityExerciseRequest, null=True, blank=True, on_delete=models.SET_NULL, related_name='payroll_tax_events')
    staff = models.ForeignKey(EntityStaff, null=True, blank=True, on_delete=models.SET_NULL, related_name='equity_payroll_tax_events')
    event_type = models.CharField(max_length=40, default='exercise_withholding')
    gross_amount = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    withholding_amount = models.DecimalField(max_digits=18, decimal_places=2, default=0)
    payroll_sync_status = models.CharField(max_length=20, choices=PayrollSyncStatus.choices, default=PayrollSyncStatus.QUEUED)
    tax_jurisdiction = models.CharField(max_length=120, blank=True, default='')
    reference_number = models.CharField(max_length=120, blank=True, default='')
    source_account = models.ForeignKey(ChartOfAccounts, null=True, blank=True, on_delete=models.SET_NULL, related_name='equity_payroll_tax_source_events')
    destination_account = models.ForeignKey(ChartOfAccounts, null=True, blank=True, on_delete=models.SET_NULL, related_name='equity_payroll_tax_destination_events')
    details = models.JSONField(default=dict, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-created_at']


class EquityExternalAdapterConfig(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_adapter_configs')
    adapter_type = models.CharField(max_length=20, choices=ExternalAdapterType.choices)
    provider_name = models.CharField(max_length=120)
    base_url = models.URLField()
    api_key = models.CharField(max_length=255, blank=True, default='')
    auth_scheme = models.CharField(max_length=40, default='Bearer')
    endpoint_path = models.CharField(max_length=255, blank=True, default='')
    webhook_secret = models.CharField(max_length=255, blank=True, default='')
    default_headers = models.JSONField(default=dict, blank=True)
    adapter_settings = models.JSONField(default=dict, blank=True)
    is_active = models.BooleanField(default=True)
    last_synced_at = models.DateTimeField(null=True, blank=True)
    last_error = models.TextField(blank=True, default='')
    created_by = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='created_equity_adapter_configs')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['adapter_type', 'provider_name']
        unique_together = ('workspace', 'adapter_type', 'provider_name')


class EquityDeliveryLog(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    workspace = models.ForeignKey(Entity, on_delete=models.CASCADE, related_name='equity_delivery_logs')
    grant = models.ForeignKey(EquityGrant, null=True, blank=True, on_delete=models.SET_NULL, related_name='delivery_logs')
    vesting_event = models.ForeignKey(EquityVestingEvent, null=True, blank=True, on_delete=models.SET_NULL, related_name='delivery_logs')
    exercise_request = models.ForeignKey(EquityExerciseRequest, null=True, blank=True, on_delete=models.SET_NULL, related_name='delivery_logs')
    certificate = models.ForeignKey(EquityShareCertificate, null=True, blank=True, on_delete=models.SET_NULL, related_name='delivery_logs')
    payroll_tax_event = models.ForeignKey(EquityPayrollTaxEvent, null=True, blank=True, on_delete=models.SET_NULL, related_name='delivery_logs')
    recipient_user = models.ForeignKey(settings.AUTH_USER_MODEL, null=True, blank=True, on_delete=models.SET_NULL, related_name='equity_delivery_logs')
    recipient_email = models.EmailField(blank=True, default='')
    channel = models.CharField(max_length=20, choices=DeliveryChannel.choices)
    event_name = models.CharField(max_length=120)
    status = models.CharField(max_length=20, choices=DeliveryStatus.choices, default=DeliveryStatus.QUEUED)
    subject = models.CharField(max_length=255, blank=True, default='')
    message = models.TextField(blank=True, default='')
    document_payload = models.JSONField(default=dict, blank=True)
    document_file = models.FileField(upload_to='equity/delivery_documents/', null=True, blank=True)
    provider_response = models.JSONField(default=dict, blank=True)
    error_message = models.TextField(blank=True, default='')
    delivered_at = models.DateTimeField(null=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-created_at']
