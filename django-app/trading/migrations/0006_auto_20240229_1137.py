# Generated by Django 3.2.24 on 2024-02-29 11:37

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('trading', '0005_alter_strategy_equity_percent'),
    ]

    operations = [
        migrations.AlterField(
            model_name='strategy',
            name='inverse',
            field=models.BooleanField(blank=True, default=False, null=True),
        ),
        migrations.AlterField(
            model_name='strategy',
            name='market_data_source',
            field=models.CharField(blank=True, default=None, max_length=255, null=True),
        ),
        migrations.AlterField(
            model_name='strategy',
            name='pyramid',
            field=models.BooleanField(blank=True, default=False, null=True),
        ),
    ]
