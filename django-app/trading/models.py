from django.db import models

class Strategy(models.Model):
    strategy_id = models.CharField(max_length=255, unique=True)
    name = models.CharField(max_length=255)
    symbol = models.CharField(max_length=255)
    initial_balance = models.DecimalField(max_digits=19, decimal_places=4)
    lots = models.IntegerField()
    equity_percent = models.DecimalField(max_digits=5, decimal_places=2)
    inverse = models.BooleanField(default=False)
    pyramid = models.BooleanField(default=False)
    market_data_source = models.CharField(max_length=255)

    def __str__(self):
        return self.name

class Trade(models.Model):
    strategy = models.ForeignKey(Strategy, on_delete=models.CASCADE)
    action = models.CharField(max_length=10)
    quantity = models.IntegerField()
    direction = models.CharField(max_length=10)
    entry_time = models.DateTimeField()
    exit_time = models.DateTimeField(null=True, blank=True)
    open_price = models.DecimalField(max_digits=19, decimal_places=4)
    close_price = models.DecimalField(max_digits=19, decimal_places=4, null=True, blank=True)
    leverage = models.DecimalField(max_digits=10, decimal_places=2, default=1.00)
    pnl = models.DecimalField(max_digits=19, decimal_places=4, null=True, blank=True)
    pnl_percentage = models.DecimalField(max_digits=10, decimal_places=2, null=True, blank=True)
    status = models.CharField(max_length=10)

    def __str__(self):
        return f"{self.strategy.name} - {self.action}"
