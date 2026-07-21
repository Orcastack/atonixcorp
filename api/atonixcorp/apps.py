from django.apps import AppConfig


class AtonixCorpConfig(AppConfig):
    default_auto_field = 'django.db.models.BigAutoField'
    name = 'atonixcorp'
    label = 'finances'

    def ready(self):
        from . import signals  # noqa: F401
