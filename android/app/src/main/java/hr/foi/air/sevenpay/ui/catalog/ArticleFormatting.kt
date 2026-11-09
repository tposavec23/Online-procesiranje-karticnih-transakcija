package hr.foi.air.sevenpay.ui.catalog

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hr.foi.air.sevenpay.ui.theme.Dostupno
import hr.foi.air.sevenpay.ui.theme.DostupnoPozadina
import hr.foi.air.sevenpay.ui.theme.MaloNaZalihi
import hr.foi.air.sevenpay.ui.theme.MaloNaZalihiPozadina
import hr.foi.air.sevenpay.ui.theme.Nedostupno
import hr.foi.air.sevenpay.ui.theme.NedostupnoPozadina
import java.math.BigDecimal
import java.util.Locale

private const val LOW_STOCK_THRESHOLD = 2

internal enum class StockStatus(val background: Color, val textColor: Color) {
    IN_STOCK(DostupnoPozadina, Dostupno),
    LOW(MaloNaZalihiPozadina, MaloNaZalihi),
    UNAVAILABLE(NedostupnoPozadina, Nedostupno)
}

internal fun stockStatus(quantity: Int) = when {
    quantity <= 0 -> StockStatus.UNAVAILABLE
    quantity <= LOW_STOCK_THRESHOLD -> StockStatus.LOW
    else -> StockStatus.IN_STOCK
}

internal fun formatPrice(cents: Long): String =
    String.format(Locale.GERMANY, "%.2f €", BigDecimal.valueOf(cents, 2))

@Composable
internal fun StockChip(quantity: Int) {
    val status = stockStatus(quantity)
    val text = when (status) {
        StockStatus.IN_STOCK -> "Na zalihi: $quantity kom"
        StockStatus.LOW -> "Malo na zalihi: $quantity kom"
        StockStatus.UNAVAILABLE -> "Nedostupno"
    }

    Box(
        modifier = Modifier
            .background(status.background, RoundedCornerShape(50))
            .padding(horizontal = 10.dp, vertical = 4.dp)
    ) {
        Text(
            text = text,
            fontSize = 12.sp,
            fontWeight = FontWeight.Medium,
            color = status.textColor
        )
    }
}