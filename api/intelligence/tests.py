from django.contrib.auth.models import User
from rest_framework import status
from rest_framework.test import APITestCase

from atonixcorp.models import Entity, Organization, PlatformAuditEvent, TaxProfile
from workspaces.services import WorkspaceService

from .models import AIInteraction, WorkspaceAIPrecedent, WorkspaceAIProfile


class WorkspaceIntelligenceAPITests(APITestCase):
    def setUp(self):
        self.owner = User.objects.create_user(
            username='ai-owner',
            email='ai-owner@example.com',
            password='password123',
        )
        self.organization = Organization.objects.create(
            owner=self.owner,
            name='AI Org',
            slug='ai-org',
            primary_country='ZA',
            primary_currency='ZAR',
        )
        self.entity = Entity.objects.create(
            organization=self.organization,
            name='ZA Trading Entity',
            country='ZA',
            entity_type='corporation',
            status='active',
            local_currency='ZAR',
        )
        TaxProfile.objects.create(
            entity=self.entity,
            country='ZA',
            tax_rules={'corporate_rate': '0.27'},
        )
        self.workspace = WorkspaceService.create_workspace(
            self.owner,
            {'name': 'AI Workspace', 'linked_entity_id': self.entity.id},
        )
        self.client.force_authenticate(self.owner)

    def test_profile_can_be_retrieved_and_updated(self):
        response = self.client.patch(
            f'/api/v1/workspaces/{self.workspace.id}/ai/profile',
            {
                'tax_preferences': {'depreciation_method': 'straight_line'},
                'risk_tolerance': 'low',
                'tone_preferences': {'style': 'concise'},
            },
            format='json',
        )

        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(response.data['tax_preferences']['depreciation_method'], 'straight_line')
        self.assertEqual(response.data['risk_tolerance'], 'low')
        self.assertTrue(WorkspaceAIProfile.objects.filter(workspace=self.workspace).exists())

    def test_tax_treatment_endpoint_returns_audited_deterministic_result(self):
        response = self.client.post(
            f'/api/v1/workspaces/{self.workspace.id}/ai/tax-treatment',
            {
                'jurisdiction': 'ZA',
                'period': '2025-Q1',
                'transaction': {
                    'type': 'expense',
                    'amount': '1000.00',
                    'currency': 'ZAR',
                    'date': '2025-02-15',
                    'description': 'Cloud hosting invoice',
                    'tags': ['opex'],
                },
                'options': {'model': 'sonnet', 'max_tokens': 1200},
            },
            format='json',
        )

        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(response.data['proposed_treatment']['category'], 'deductible_expense')
        self.assertEqual(response.data['deterministic_calculation']['tax_effect_amount'], '270.00')
        self.assertEqual(AIInteraction.objects.filter(workspace=self.workspace, intent='tax_treatment_explanation').count(), 1)
        self.assertTrue(
            PlatformAuditEvent.objects.filter(
                domain='ai',
                event_type='ai.tax_treatment_explained',
                workspace_id=self.workspace.id,
            ).exists()
        )

    def test_feedback_creates_workspace_precedent(self):
        explain_response = self.client.post(
            f'/api/v1/workspaces/{self.workspace.id}/ai/tax-treatment',
            {
                'jurisdiction': 'ZA',
                'period': '2025-Q1',
                'transaction': {
                    'type': 'expense',
                    'amount': '1000.00',
                    'currency': 'ZAR',
                    'date': '2025-02-15',
                    'description': 'Cloud hosting invoice',
                    'tags': ['opex'],
                },
            },
            format='json',
        )
        interaction_id = explain_response.data['interaction_id']

        feedback_response = self.client.post(
            f'/api/v1/workspaces/{self.workspace.id}/ai/interactions/{interaction_id}/feedback',
            {
                'outcome': 'accepted',
                'reason': 'Matches our recurring treatment for hosting costs.',
            },
            format='json',
        )

        self.assertEqual(feedback_response.status_code, status.HTTP_200_OK)
        self.assertEqual(feedback_response.data['feedback_outcome'], 'accepted')
        self.assertTrue(WorkspaceAIPrecedent.objects.filter(workspace=self.workspace, feedback_outcome='accepted').exists())
        profile = WorkspaceAIProfile.objects.get(workspace=self.workspace)
        self.assertEqual(profile.feedback_history[-1]['outcome'], 'accepted')
