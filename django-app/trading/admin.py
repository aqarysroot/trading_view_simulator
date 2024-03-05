from django.contrib import admin
from .models import Strategy, Trade

class StrategyAdmin(admin.ModelAdmin):
    list_display = ('name', 'symbol', 'initial_balance', 'lots', 'equity_percent', 'inverse', 'pyramid', 'market_data_source')
    search_fields = ('name', 'symbol')

class TradeAdmin(admin.ModelAdmin):
    list_display = ('strategy', 'action', 'quantity', 'direction', 'entry_time', 'exit_time', 'open_price', 'close_price', 'pnl', 'pnl_percentage')
    list_filter = ('action', 'direction', 'entry_time', 'exit_time')
    search_fields = ('strategy__name', 'action')

admin.site.register(Strategy, StrategyAdmin)
admin.site.register(Trade, TradeAdmin)
