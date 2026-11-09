package hr.foi.air.sevenpay.ui.catalogue.detail

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import hr.foi.air.sevenpay.ui.theme.CreamBackground
import hr.foi.air.sevenpay.ui.theme.Raspberry
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme
import hr.foi.air.sevenpay.ui.theme.SecondaryText

@Composable
fun ArticleDetailMessage(
    message: String,
    onBack: () -> Unit,
    modifier: Modifier = Modifier,
    onRetry: (() -> Unit)? = null,
    isError: Boolean = false
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .background(CreamBackground)
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = message,
            color = if (isError) MaterialTheme.colorScheme.error else SecondaryText,
            textAlign = TextAlign.Center
        )
        Spacer(Modifier.height(12.dp))

        if (onRetry != null) {
            Button(
                onClick = onRetry,
                colors = ButtonDefaults.buttonColors(containerColor = Raspberry)
            ) {
                Text("Pokušaj ponovno")
            }
            Spacer(Modifier.height(8.dp))
            OutlinedButton(onClick = onBack) {
                Text("Natrag na popis")
            }
        } else {
            Button(
                onClick = onBack,
                colors = ButtonDefaults.buttonColors(containerColor = Raspberry)
            ) {
                Text("Natrag na popis")
            }
        }
    }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun ArticleDetailMessagePreview() {
    SevenPayTheme {
        ArticleDetailMessage(message = "Artikl više ne postoji.", onBack = {})
    }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun ArticleDetailErrorPreview() {
    SevenPayTheme {
        ArticleDetailMessage(
            message = "Dohvat detalja artikla nije uspio. Provjerite vezu i pokušajte ponovno.",
            onBack = {},
            onRetry = {},
            isError = true
        )
    }
}