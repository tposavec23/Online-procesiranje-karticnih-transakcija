package hr.foi.air.sevenpay.ui.catalogue

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.ui.theme.CreamBackground
import hr.foi.air.sevenpay.ui.theme.Raspberry
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme
import hr.foi.air.sevenpay.ui.theme.SecondaryText

@Composable
fun CatalogueScreen(
    state: CatalogueUiState,
    onArticleClick: (Article) -> Unit,
    onRetry: () -> Unit,
    onLoadMore: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .background(CreamBackground)
            .padding(horizontal = 16.dp)
    ) {
        Spacer(Modifier.height(16.dp))
        Text("Artikli", fontSize = 28.sp, fontWeight = FontWeight.Bold)
        Spacer(Modifier.height(16.dp))

        when (state) {
            CatalogueUiState.Loading -> CenteredContent {
                CircularProgressIndicator(color = Raspberry)
            }

            CatalogueUiState.Empty -> CenteredContent {
                Text(
                    "Trenutno nema dostupnih artikala.",
                    color = SecondaryText,
                    textAlign = TextAlign.Center
                )
            }

            is CatalogueUiState.Error -> CenteredContent {
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Text(
                        state.message,
                        color = MaterialTheme.colorScheme.error,
                        textAlign = TextAlign.Center
                    )
                    Spacer(Modifier.height(12.dp))
                    Button(
                        onClick = onRetry,
                        colors = ButtonDefaults.buttonColors(containerColor = Raspberry)
                    ) {
                        Text("Pokušaj ponovno")
                    }
                }
            }

            is CatalogueUiState.Success -> {
                Text(
                    "Prikazano artikala: ${state.articles.size}",
                    fontSize = 14.sp,
                    color = SecondaryText
                )
                Spacer(Modifier.height(8.dp))
                LazyColumn(
                    verticalArrangement = Arrangement.spacedBy(12.dp),
                    contentPadding = PaddingValues(bottom = 16.dp)
                ) {
                    items(state.articles, key = { it.id }) { article ->
                        ArticleCard(
                            article = article,
                            onClick = { onArticleClick(article) }
                        )
                    }
                    if (state.hasMore) {
                        item {
                            LoadMoreFooter(
                                isLoading = state.isLoadingMore,
                                failed = state.loadMoreFailed,
                                onClick = onLoadMore
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun LoadMoreFooter(
    isLoading: Boolean,
    failed: Boolean,
    onClick: () -> Unit
) {
    Column(
        modifier = Modifier.fillMaxWidth(),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        if (failed) {
            Text(
                "Dohvat sljedećih artikala nije uspio.",
                color = MaterialTheme.colorScheme.error,
                textAlign = TextAlign.Center
            )
            Spacer(Modifier.height(8.dp))
        }
        if (isLoading) {
            CircularProgressIndicator(color = Raspberry)
        } else {
            Button(
                onClick = onClick,
                colors = ButtonDefaults.buttonColors(containerColor = Raspberry)
            ) {
                Text("Prikaži više")
            }
        }
    }
}

@Composable
private fun CenteredContent(content: @Composable () -> Unit) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        contentAlignment = Alignment.Center
    ) { content() }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun CatalogueScreenPreview() {
    SevenPayTheme {
        CatalogueScreen(CatalogueUiState.Success(sampleArticles, hasMore = true), {}, {}, {})
    }
}

@Preview(showBackground = true)
@Composable
private fun CatalogueScreenEmptyPreview() {
    SevenPayTheme {
        CatalogueScreen(CatalogueUiState.Empty, {}, {}, {})
    }
}

@Preview(showBackground = true)
@Composable
private fun CatalogueScreenErrorPreview() {
    SevenPayTheme {
        CatalogueScreen(CatalogueUiState.Error("Dohvat artikala nije uspio."), {}, {}, {})
    }
}