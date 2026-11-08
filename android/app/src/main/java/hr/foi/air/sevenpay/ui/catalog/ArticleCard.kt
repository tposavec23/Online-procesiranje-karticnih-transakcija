package hr.foi.air.sevenpay.ui.catalog

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.ui.theme.InStockText
import hr.foi.air.sevenpay.ui.theme.InStockBackground
import hr.foi.air.sevenpay.ui.theme.LowStockText
import hr.foi.air.sevenpay.ui.theme.LowStockBackground
import hr.foi.air.sevenpay.ui.theme.UnavailableText
import hr.foi.air.sevenpay.ui.theme.UnavailableBackground
import hr.foi.air.sevenpay.ui.theme.Border
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme
import hr.foi.air.sevenpay.ui.theme.SecondaryText
import java.util.Locale
import java.math.BigDecimal

private const val LOW_STOCK_THRESHOLD = 2

private enum class StockStatus(val background: Color, val textColor: Color) {
    IN_STOCK(InStockBackground, InStockText),
    LOW(LowStockBackground, LowStockText),
    UNAVAILABLE(UnavailableBackground, UnavailableText)
}

private fun stockStatus(quantity: Int) = when {
    quantity <= 0 -> StockStatus.UNAVAILABLE
    quantity <= LOW_STOCK_THRESHOLD -> StockStatus.LOW
    else -> StockStatus.IN_STOCK
}

@Composable
fun ArticleCard(
    article: Article,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    val status = stockStatus(article.availableQuantity)

    Card(
        onClick = onClick,
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = BorderStroke(1.dp, Border)
    ) {
        Column(Modifier.padding(16.dp)) {
            Text(
                text = article.code,
                fontSize = 12.sp,
                fontWeight = FontWeight.SemiBold,
                color = SecondaryText
            )
            Spacer(Modifier.height(4.dp))
            Text(
                text = article.name,
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold
            )
            Text(
                text = article.description,
                fontSize = 14.sp,
                color = SecondaryText
            )
            Spacer(Modifier.height(12.dp))
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(
                    text = String.format(Locale.GERMANY, "%.2f €", BigDecimal.valueOf(article.price, 2)),
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold
                )
                Spacer(Modifier.width(12.dp))
                StockChip(
                    text = when (status) {
                        StockStatus.IN_STOCK -> "Na zalihi: ${article.availableQuantity} kom"
                        StockStatus.LOW -> "Malo na zalihi: ${article.availableQuantity} kom"
                        StockStatus.UNAVAILABLE -> "Nedostupno"
                    },
                    status = status
                )
            }
        }
    }
}

@Composable
private fun StockChip(text: String, status: StockStatus) {
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

@Preview(showBackground = true)
@Composable
private fun ArticleCardPreview() {
    SevenPayTheme {
        ArticleCard(article = sampleArticles.first(), onClick = {})
    }
}