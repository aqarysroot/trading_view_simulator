from django.db import models


class Strategy(models.Model):
    strategy_id = models.CharField(max_length=255)
    symbol = models.CharField(max_length=10)
    initial_balance = models.FloatField()
    lots = models.IntegerField()
    equity_percent = models.FloatField()
    inverse = models.BooleanField(default=False)
    pyramid = models.BooleanField(default=False)
    market_data_source = models.CharField(max_length=255)
    open_price = models.DecimalField(max_digits=10, decimal_places=2, default=None, blank=True, null=True)
    close_price = models.DecimalField(max_digits=10, decimal_places=2, default=None, blank=True, null=True)
    entry_time = models.DateTimeField(default=None, blank=True, null=True)
    exit_time = models.DateTimeField(default=None, blank=True, null=True)
    pnl = models.DecimalField(max_digits=10, decimal_places=2, verbose_name="Profit and Loss (USD)", default=None, blank=True, null=True)

    def __str__(self):
        return self.name
    