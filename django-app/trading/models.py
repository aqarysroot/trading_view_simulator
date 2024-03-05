from django.db import models
from django.utils import timezone

class Strategy(models.Model):
    strategy_id = models.CharField(max_key=255, unique=True)
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
    action = models.CharField(max_length=10)  # 'buy' or 'sell'
    quantity = models.IntegerField()
    direction = models.CharField(max_length=10)  # 'Long' or 'Short'
    buy_datetime = models.DateTimeField(null=True, blank=True)
    sell_datetime = models.DateTimeField(null=True, blank=True)
    buy_price = models.DecimalField(max_digits=19, decimal_places=4, null=True, blank=True)
    sell_price = models.DecimalField(max_digits=19, decimal_places=4, null=True, blank=True)
    leverage = models.DecimalField(max_digits=10, decimal_places=2, default=1.00)
    pnl = models.DecimalField(max_digits=19, decimal_places=4, null=True, blank=True)
    pnl_percentage = models.DecimalField(max_digits=10, decimal_places=2, null=True, blank=True)

    def __str__(self):
        return f"{self.strategy.name} - {self.action}"

    def save(self, *args, **kwargs):
        if self.buy_price and self.sell_price:
            self.pnl = (self.sell_price - self.buy_price) * self.quantity
            self.pnl_percentage = (self.pnl / (self.buy_price * self.quantity)) * 100
        super().save(*args, **kwargs)
