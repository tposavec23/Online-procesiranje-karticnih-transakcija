package hr.foi.air.sevenpay.ui.catalogue.detail

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.BakeryDining
import androidx.compose.material.icons.filled.Cake
import androidx.compose.material.icons.filled.Inventory2
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.ui.catalogue.StockChip
import hr.foi.air.sevenpay.ui.catalogue.formatPrice
import hr.foi.air.sevenpay.ui.catalogue.sampleArticles
import hr.foi.air.sevenpay.ui.theme.HeroPeach
import hr.foi.air.sevenpay.ui.theme.CreamBackground
import hr.foi.air.sevenpay.ui.theme.Raspberry
import hr.foi.air.sevenpay.ui.theme.PrimaryText
import hr.foi.air.sevenpay.ui.theme.Border
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme
import hr.foi.air.sevenpay.ui.theme.SecondaryText
import hr.foi.air.sevenpay.ui.theme.LightRose
import hr.foi.air.sevenpay.ui.theme.DarkRaspberry

private fun categoryIcon(category: String): ImageVector = when (category) {
    "Torte" -> Icons.Filled.Cake
    "Kolači" -> Icons.Filled.BakeryDining
    else -> Icons.Filled.Inventory2
}

@Composable
fun ArticleDetailScreen(
    article: Article,
    onBack: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .background(CreamBackground)
            .verticalScroll(rememberScrollState())
            .padding(horizontal = 16.dp)
    ) {
        Spacer(Modifier.height(16.dp))
        Row(verticalAlignment = Alignment.CenterVertically) {
            IconButton(
                onClick = onBack,
                modifier = Modifier
                    .size(44.dp)
                    .background(Color.White, RoundedCornerShape(12.dp))
                    .border(1.dp, Border, RoundedCornerShape(12.dp))
            ) {
                Icon(
                    imageVector = Icons.AutoMirrored.Filled.ArrowBack,
                    contentDescription = "Natrag na popis artikala"
                )
            }
            Spacer(Modifier.width(12.dp))
            Text("Detalji artikla", fontSize = 20.sp, fontWeight = FontWeight.Bold)
        }
        Spacer(Modifier.height(16.dp))

        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(180.dp)
                .background(HeroPeach, RoundedCornerShape(24.dp)),
            contentAlignment = Alignment.Center
        ) {
            Icon(
                imageVector = categoryIcon(article.category),
                contentDescription = article.category,
                tint = Raspberry,
                modifier = Modifier.size(72.dp)
            )
        }
        Spacer(Modifier.height(16.dp))

        Row(verticalAlignment = Alignment.CenterVertically) {
            Text(
                text = article.code,
                fontSize = 12.sp,
                fontWeight = FontWeight.SemiBold,
                color = SecondaryText
            )
            Spacer(Modifier.width(8.dp))
            CategoryTag(article.category)
        }
        Spacer(Modifier.height(4.dp))
        Text(article.name, fontSize = 28.sp, fontWeight = FontWeight.Bold)
        Spacer(Modifier.height(8.dp))
        Row(verticalAlignment = Alignment.CenterVertically) {
            Text(
                text = formatPrice(article.price),
                fontSize = 28.sp,
                fontWeight = FontWeight.Bold
            )
            Spacer(Modifier.width(12.dp))
            StockChip(article.availableQuantity)
        }
        Spacer(Modifier.height(16.dp))

        Card(
            modifier = Modifier.fillMaxWidth(),
            shape = RoundedCornerShape(16.dp),
            colors = CardDefaults.cardColors(containerColor = Color.White),
            border = BorderStroke(1.dp, Border)
        ) {
            Column(Modifier.padding(16.dp)) {
                Text("Opis", fontSize = 14.sp, fontWeight = FontWeight.Bold)
                Spacer(Modifier.height(4.dp))
                Text(article.description, fontSize = 14.sp, color = PrimaryText)
            }
        }
        Spacer(Modifier.height(16.dp))
    }
}

@Composable
private fun CategoryTag(text: String) {
    Box(
        modifier = Modifier
            .background(LightRose, RoundedCornerShape(50))
            .padding(horizontal = 10.dp, vertical = 4.dp)
    ) {
        Text(
            text = text,
            fontSize = 12.sp,
            fontWeight = FontWeight.Medium,
            color = DarkRaspberry
        )
    }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun ArticleDetailScreenPreview() {
    SevenPayTheme {
        ArticleDetailScreen(sampleArticles.first(), onBack = {})
    }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun ArticleDetailScreenUnavailablePreview() {
    SevenPayTheme {
        ArticleDetailScreen(sampleArticles[2], onBack = {})
    }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun ArticleDetailScreenLowStockPreview() {
    SevenPayTheme {
        ArticleDetailScreen(sampleArticles[3], onBack = {})
    }
}