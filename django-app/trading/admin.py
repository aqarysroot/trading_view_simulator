from django.contrib import admin
from .models import Strategy

@admin.register(Strategy)
class StrategyAdmin(admin.ModelAdmin):
    list_display = ('strategy_id', 'symbol', 'initial_balance', 'open_price', 'close_price', 'entry_time', 'exit_time', 'pnl')
    list_filter = ('strategy_id',) 
    