import json
from dataclasses import dataclass
from decimal import Decimal, InvalidOperation, ROUND_HALF_UP
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen

from django.conf import settings
from django.db import transaction
from django.db.models import Q

from atonixcorp.platform_foundation import log_platform_audit_event
from atonixcorp.models import TaxProfile
from workspaces.models import Workspace
from workspaces.services import LogService

from .models import (
    AIFeedbackOutcome,
    AIInteraction,
    AIInteractionStatus,
    GlobalKnowledgeChunk,
    WorkspaceAIPrecedent,
    WorkspaceAIProfile,
)
from .prompt_library import build_tax_treatment_prompt


DEFAULT_TAX_RATES = {
    'ZA': Decimal('0.27'),
    'US': Decimal('0.21'),
    'UK': Decimal('0.25'),
    'GLOBAL': Decimal('0.30'),
}
CAPITALIZATION_KEYWORDS = {'capital', 'capex', 'fixed_asset', 'asset', 'equipment', 'intangible'}
EQUITY_REVIEW_TYPES = {'equity_grant', 'grant', 'vesting', 'exercise', 'convertible', 'safe'}


def _quantize(amount):
    return amount.quantize(Decimal('0.01'), rounding=ROUND_HALF_UP)


def _decimal(value, default='0'):
    try:
        return Decimal(str(value if value is not None else default))
    except (InvalidOperation, TypeError, ValueError):
        return Decimal(str(default))


@dataclass
class ModelSelection:
    key: str
    model_name: str


class ClaudeModelRouter:
    @staticmethod
    def select(intent, payload, options=None):
        options = options or {}
        requested_model = (options.get('model') or '').strip().lower()
        if requested_model == 'opus':
            return ModelSelection('opus', settings.CLAUDE_OPUS_MODEL)
        if requested_model == 'sonnet':
            return ModelSelection('sonnet', settings.CLAUDE_SONNET_MODEL)

        transaction = payload.get('transaction') or {}
        serialized_size = len(json.dumps(payload, sort_keys=True))
        complexity = 0
        if serialized_size > 2500:
            complexity += 2
        if len(transaction.get('tags') or []) > 3:
            complexity += 1
        if str(transaction.get('type', '')).lower() in EQUITY_REVIEW_TYPES:
            complexity += 2
        if isinstance(payload.get('documents'), list) and len(payload['documents']) > 1:
            complexity += 2

        if complexity >= 3:
            return ModelSelection('opus', settings.CLAUDE_OPUS_MODEL)
        return ModelSelection('sonnet', settings.CLAUDE_SONNET_MODEL)


class ClaudeGateway:
    @staticmethod
    def _request_payload(*, model_name, prompt, max_tokens):
        return {
            'model': model_name,
            'max_tokens': max_tokens,
            'system': prompt['system'],
            'messages': [
                {
                    'role': 'user',
                    'content': prompt['user'],
                }
            ],
        }

    @classmethod
    def _call_api(cls, *, model_name, prompt, max_tokens):
        body = cls._request_payload(model_name=model_name, prompt=prompt, max_tokens=max_tokens)
        url = f"{settings.ANTHROPIC_API_BASE_URL.rstrip('/')}/v1/messages"
        request = Request(
            url,
            data=json.dumps(body).encode('utf-8'),
            headers={
                'content-type': 'application/json',
                'x-api-key': settings.ANTHROPIC_API_KEY,
                'anthropic-version': '2023-06-01',
            },
            method='POST',
        )
        with urlopen(request, timeout=30) as response:
            return json.loads(response.read().decode('utf-8'))

    @staticmethod
    def _fallback_response(deterministic_assessment, references, error_message=''):
        proposed_treatment = deterministic_assessment['proposed_treatment']
        calculation = deterministic_assessment['calculation']
        trace = deterministic_assessment['trace']
        answer = (
            f"Deterministic assessment suggests `{proposed_treatment['category']}` for this transaction. "
            f"Verified tax impact is {calculation['tax_effect_amount']} {calculation['currency']} "
            f"with direction `{calculation['impact_direction']}`."
        )
        if error_message:
            answer = f'{answer} External model unavailable: {error_message}.'
        return {
            'answer': answer,
            'proposed_treatment': proposed_treatment,
            'references': references,
            'confidence': deterministic_assessment['confidence'],
            'requires_human_review': deterministic_assessment['requires_human_review'],
            'trace': trace,
        }

    @classmethod
    def generate_response(cls, *, model_name, prompt, deterministic_assessment, references, max_tokens):
        should_call_model = settings.AI_ENABLE_EXTERNAL_MODELS and bool(settings.ANTHROPIC_API_KEY)
        if not should_call_model:
            return cls._fallback_response(deterministic_assessment, references), {}, AIInteractionStatus.FALLBACK

        try:
            raw = cls._call_api(model_name=model_name, prompt=prompt, max_tokens=max_tokens)
            content_blocks = raw.get('content') or []
            text = '\n'.join(block.get('text', '') for block in content_blocks if isinstance(block, dict))
            try:
                parsed = json.loads(text)
            except json.JSONDecodeError:
                parsed = {
                    'answer': text or 'Model response did not provide structured JSON.',
                    'proposed_treatment': deterministic_assessment['proposed_treatment'],
                    'references': references,
                    'confidence': deterministic_assessment['confidence'],
                    'requires_human_review': deterministic_assessment['requires_human_review'],
                }
            parsed.setdefault('proposed_treatment', deterministic_assessment['proposed_treatment'])
            parsed.setdefault('references', references)
            parsed.setdefault('confidence', deterministic_assessment['confidence'])
            parsed.setdefault('requires_human_review', deterministic_assessment['requires_human_review'])
            return parsed, raw, AIInteractionStatus.SUCCEEDED
        except (HTTPError, URLError, TimeoutError, ValueError) as exc:
            return cls._fallback_response(deterministic_assessment, references, str(exc)), {'error': str(exc)}, AIInteractionStatus.FALLBACK


class DeterministicTaxEngine:
    @staticmethod
    def _rate_from_profile(workspace, jurisdiction, transaction_type):
        linked_entity = workspace.linked_entity
        rate = None
        rate_source = 'default_jurisdiction_rate'
        if linked_entity:
            profile = TaxProfile.objects.filter(entity=linked_entity, country__iexact=jurisdiction).first()
            if profile:
                rules = profile.tax_rules or {}
                type_key_map = {
                    'vat': 'vat_rate',
                    'withholding': 'withholding_rate',
                }
                preferred_key = type_key_map.get(transaction_type, 'corporate_rate')
                raw_rate = rules.get(preferred_key, rules.get('default_rate'))
                if raw_rate not in (None, ''):
                    rate = _decimal(raw_rate)
                    rate_source = f'tax_profile:{linked_entity.id}:{jurisdiction}'
        if rate is None or rate <= 0:
            rate = DEFAULT_TAX_RATES.get(jurisdiction.upper(), DEFAULT_TAX_RATES['GLOBAL'])
        if rate > 1:
            rate = rate / Decimal('100')
        return rate, rate_source

    @classmethod
    def assess_transaction(cls, *, workspace, jurisdiction, payload):
        transaction = payload['transaction']
        amount = _quantize(_decimal(transaction.get('amount')))
        currency = (transaction.get('currency') or 'USD').upper()
        transaction_type = str(transaction.get('type') or 'expense').lower()
        tags = {str(tag).strip().lower() for tag in (transaction.get('tags') or [])}
        raw_data = transaction.get('raw_data') if isinstance(transaction.get('raw_data'), dict) else {}
        memo = ' '.join(
            str(item) for item in [
                transaction.get('description', ''),
                raw_data.get('description', ''),
                raw_data.get('memo', ''),
            ] if item
        ).lower()
        rate, rate_source = cls._rate_from_profile(workspace, jurisdiction, transaction_type)

        trace = [
            {'step': 'normalize_transaction', 'result': {'type': transaction_type, 'amount': str(amount), 'currency': currency}},
            {'step': 'resolve_tax_rate', 'result': {'rate': str(rate), 'source': rate_source}},
        ]
        requires_human_review = False
        confidence = 'medium'

        is_capitalized = bool(tags & CAPITALIZATION_KEYWORDS) or any(keyword in memo for keyword in ['capex', 'capital', 'fixed asset', 'equipment purchase'])
        if transaction_type in {'expense', 'vendor_bill', 'bill', 'cost'}:
            if is_capitalized:
                trace.append({'step': 'classify_expense', 'result': 'capitalized'})
                proposed_treatment = {
                    'category': 'capital_expenditure',
                    'rationale': 'Capitalization indicators were detected in tags or transaction memo, so immediate deduction is blocked in the deterministic engine.',
                }
                calculation = {
                    'tax_rate': str(rate),
                    'tax_base_amount': '0.00',
                    'tax_effect_amount': '0.00',
                    'impact_direction': 'deferred',
                    'currency': currency,
                }
            else:
                trace.append({'step': 'classify_expense', 'result': 'deductible'})
                tax_effect = _quantize(amount * rate)
                proposed_treatment = {
                    'category': 'deductible_expense',
                    'rationale': 'Default expense treatment is deductible unless capitalization indicators are present.',
                }
                calculation = {
                    'tax_rate': str(rate),
                    'tax_base_amount': str(amount),
                    'tax_effect_amount': str(tax_effect),
                    'impact_direction': 'reduction',
                    'currency': currency,
                }
        elif transaction_type in {'sale', 'income', 'revenue', 'invoice', 'dividend'}:
            trace.append({'step': 'classify_income', 'result': 'taxable'})
            tax_effect = _quantize(amount * rate)
            proposed_treatment = {
                'category': 'taxable_income',
                'rationale': 'Revenue-like transaction types are treated as taxable income in the deterministic engine.',
            }
            calculation = {
                'tax_rate': str(rate),
                'tax_base_amount': str(amount),
                'tax_effect_amount': str(tax_effect),
                'impact_direction': 'increase',
                'currency': currency,
            }
        elif transaction_type in EQUITY_REVIEW_TYPES:
            trace.append({'step': 'classify_equity_event', 'result': 'manual_review_required'})
            requires_human_review = True
            confidence = 'low'
            proposed_treatment = {
                'category': 'equity_event_manual_review',
                'rationale': 'Equity events require jurisdiction-specific rules and grant details before final tax treatment can be confirmed.',
            }
            calculation = {
                'tax_rate': str(rate),
                'tax_base_amount': '0.00',
                'tax_effect_amount': '0.00',
                'impact_direction': 'review_required',
                'currency': currency,
            }
        else:
            trace.append({'step': 'classify_unknown', 'result': 'manual_review_required'})
            requires_human_review = True
            confidence = 'low'
            proposed_treatment = {
                'category': 'manual_review_required',
                'rationale': 'The transaction type is outside the deterministic Phase 1 ruleset.',
            }
            calculation = {
                'tax_rate': str(rate),
                'tax_base_amount': '0.00',
                'tax_effect_amount': '0.00',
                'impact_direction': 'review_required',
                'currency': currency,
            }

        return {
            'proposed_treatment': proposed_treatment,
            'calculation': calculation,
            'trace': trace,
            'confidence': confidence,
            'requires_human_review': requires_human_review,
        }


class WorkspaceIntelligenceService:
    @staticmethod
    def get_or_create_profile(workspace):
        return WorkspaceAIProfile.objects.get_or_create(workspace=workspace)[0]

    @staticmethod
    def workspace_summary(workspace, profile):
        linked_entity = workspace.linked_entity
        organization = linked_entity.organization if linked_entity else None
        return {
            'workspace_id': str(workspace.id),
            'workspace_name': workspace.name,
            'region': getattr(organization, 'primary_country', '') if organization else '',
            'linked_entity_id': linked_entity.id if linked_entity else None,
            'linked_entity_name': linked_entity.name if linked_entity else '',
            'organization_id': organization.id if organization else None,
            'organization_name': organization.name if organization else '',
            'tax_preferences': profile.tax_preferences,
            'equity_policies': profile.equity_policies,
            'compliance_profile': profile.compliance_profile,
            'risk_tolerance': profile.risk_tolerance,
            'language_preferences': profile.language_preferences,
            'tone_preferences': profile.tone_preferences,
        }

    @staticmethod
    def global_context(jurisdiction, topic='tax', limit=5):
        chunks = GlobalKnowledgeChunk.objects.filter(
            Q(jurisdiction__iexact=jurisdiction) | Q(jurisdiction='GLOBAL'),
            topic__iexact=topic,
        ).order_by('-effective_date', 'chunk_index')[:limit]
        return [
            {
                'id': str(chunk.id),
                'text': chunk.text,
                'jurisdiction': chunk.jurisdiction,
                'topic': chunk.topic,
                'source': chunk.source,
                'effective_date': chunk.effective_date.isoformat() if chunk.effective_date else None,
            }
            for chunk in chunks
        ]

    @staticmethod
    def workspace_context(workspace, intent, limit=5):
        precedents = WorkspaceAIPrecedent.objects.filter(
            workspace=workspace,
            intent=intent,
            active=True,
            feedback_outcome__in=[AIFeedbackOutcome.ACCEPTED, AIFeedbackOutcome.OVERRIDDEN],
        ).order_by('-updated_at')[:limit]
        return [
            {
                'id': str(precedent.id),
                'feedback_outcome': precedent.feedback_outcome,
                'example_payload': precedent.example_payload,
                'final_treatment': precedent.final_treatment,
                'rationale': precedent.rationale,
            }
            for precedent in precedents
        ]


class AIOrchestrationService:
    TAX_TREATMENT_INTENT = 'tax_treatment_explanation'

    @classmethod
    @transaction.atomic
    def explain_tax_treatment(cls, *, workspace, user, payload):
        profile = WorkspaceIntelligenceService.get_or_create_profile(workspace)
        options = payload.get('options') or {}
        model_selection = ClaudeModelRouter.select(cls.TAX_TREATMENT_INTENT, payload, options)
        deterministic_assessment = DeterministicTaxEngine.assess_transaction(
            workspace=workspace,
            jurisdiction=payload['jurisdiction'],
            payload=payload,
        )
        global_context = WorkspaceIntelligenceService.global_context(payload['jurisdiction'])
        workspace_context = WorkspaceIntelligenceService.workspace_context(workspace, cls.TAX_TREATMENT_INTENT)
        workspace_summary = WorkspaceIntelligenceService.workspace_summary(workspace, profile)
        prompt = build_tax_treatment_prompt(
            model_name=model_selection.model_name,
            workspace_summary=workspace_summary,
            global_context=global_context,
            workspace_context=workspace_context,
            deterministic_assessment=deterministic_assessment,
            request_payload=payload,
        )

        references = [
            {
                'type': 'deterministic_rule',
                'source': step['step'],
                'snippet': json.dumps(step['result'], sort_keys=True),
            }
            for step in deterministic_assessment['trace']
        ]
        references.extend(
            {
                'type': 'global_rule',
                'source': item['source'],
                'snippet': item['text'][:280],
            }
            for item in global_context
        )
        references.extend(
            {
                'type': 'workspace_precedent',
                'source': item['id'],
                'snippet': item['rationale'][:280],
            }
            for item in workspace_context
        )

        generated_response, raw_claude_response, interaction_status = ClaudeGateway.generate_response(
            model_name=model_selection.model_name,
            prompt=prompt,
            deterministic_assessment=deterministic_assessment,
            references=references,
            max_tokens=int(options.get('max_tokens') or settings.AI_DEFAULT_MAX_TOKENS),
        )

        linked_entity = workspace.linked_entity
        organization = linked_entity.organization if linked_entity else None
        interaction = AIInteraction.objects.create(
            workspace=workspace,
            organization=organization,
            entity=linked_entity,
            user=user,
            intent=cls.TAX_TREATMENT_INTENT,
            input_payload=payload,
            claude_request=prompt,
            claude_response=generated_response,
            raw_claude_response=raw_claude_response,
            tools_used=['deterministic_tax_engine', 'global_context_lookup', 'workspace_precedent_lookup'],
            output_payload={
                'deterministic_assessment': deterministic_assessment,
                'references': references,
            },
            model_name=model_selection.model_name,
            prompt_version=prompt['version'],
            status=interaction_status,
            confidence=generated_response.get('confidence', deterministic_assessment['confidence']),
        )

        log_platform_audit_event(
            organization=organization,
            entity=linked_entity,
            workspace_id=workspace.id,
            actor=user,
            domain='ai',
            event_type='ai.tax_treatment_explained',
            resource_type='AIInteraction',
            resource_id=str(interaction.id),
            resource_name=workspace.name,
            subject_type='workspace',
            subject_id=str(workspace.id),
            action='ai.tax_treatment_explained',
            summary='AI tax treatment explanation generated',
            metadata={
                'intent': cls.TAX_TREATMENT_INTENT,
                'model_name': model_selection.model_name,
                'prompt_version': prompt['version'],
                'status': interaction_status,
            },
            context={
                'jurisdiction': payload['jurisdiction'],
                'period': payload.get('period', ''),
                'transaction_type': str((payload.get('transaction') or {}).get('type', '')),
            },
        )
        LogService.log(workspace.id, user, 'ai.tax_treatment_explained', {'interaction_id': str(interaction.id)})

        return {
            'interaction_id': str(interaction.id),
            'answer': generated_response.get('answer', ''),
            'proposed_treatment': generated_response.get('proposed_treatment', deterministic_assessment['proposed_treatment']),
            'deterministic_calculation': deterministic_assessment['calculation'],
            'deterministic_trace': deterministic_assessment['trace'],
            'references': generated_response.get('references', references),
            'confidence': generated_response.get('confidence', deterministic_assessment['confidence']),
            'requires_human_review': generated_response.get('requires_human_review', deterministic_assessment['requires_human_review']),
            'model': model_selection.model_name,
            'prompt_version': prompt['version'],
            'raw_claude_response': raw_claude_response,
        }

    @classmethod
    @transaction.atomic
    def capture_feedback(cls, *, interaction, outcome, override_treatment=None, reason=''):
        interaction.feedback_outcome = outcome
        interaction.feedback_payload = {
            'override_treatment': override_treatment or {},
            'reason': reason,
        }
        interaction.save(update_fields=['feedback_outcome', 'feedback_payload', 'updated_at'])

        final_treatment = override_treatment or interaction.claude_response.get('proposed_treatment') or {}
        precedent = WorkspaceAIPrecedent.objects.create(
            workspace=interaction.workspace,
            interaction=interaction,
            intent=interaction.intent,
            feedback_outcome=outcome,
            example_payload=interaction.input_payload,
            proposed_treatment=interaction.claude_response.get('proposed_treatment') or {},
            final_treatment=final_treatment,
            rationale=reason,
            metadata={
                'confidence': interaction.confidence,
                'model_name': interaction.model_name,
            },
            active=outcome in [AIFeedbackOutcome.ACCEPTED, AIFeedbackOutcome.OVERRIDDEN],
        )

        profile = WorkspaceIntelligenceService.get_or_create_profile(interaction.workspace)
        feedback_history = list(profile.feedback_history)
        feedback_history.append(
            {
                'interaction_id': str(interaction.id),
                'outcome': outcome,
                'reason': reason,
            }
        )
        profile.feedback_history = feedback_history[-100:]
        profile.save(update_fields=['feedback_history', 'updated_at'])

        log_platform_audit_event(
            organization=interaction.organization,
            entity=interaction.entity,
            workspace_id=interaction.workspace.id,
            actor=interaction.user,
            domain='ai',
            event_type='ai.feedback_captured',
            resource_type='AIInteraction',
            resource_id=str(interaction.id),
            resource_name=interaction.intent,
            subject_type='workspace',
            subject_id=str(interaction.workspace.id),
            action='ai.feedback_captured',
            summary='AI interaction feedback captured',
            metadata={
                'feedback_outcome': outcome,
                'precedent_id': str(precedent.id),
            },
        )
        LogService.log(interaction.workspace.id, interaction.user, 'ai.feedback_captured', {'interaction_id': str(interaction.id), 'outcome': outcome})
        return precedent


def get_workspace_or_raise(workspace_id):
    return Workspace.objects.select_related('linked_entity__organization').get(pk=workspace_id)
