package hr.foi.air.sevenpay

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import hr.foi.air.sevenpay.navigation.BottomNavigationBar
import hr.foi.air.sevenpay.ui.catalog.CatalogScreen
import hr.foi.air.sevenpay.ui.catalog.CatalogUiState
import hr.foi.air.sevenpay.ui.catalog.sampleArticles
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContent {
            SevenPayTheme {
                var selectedItem by rememberSaveable {
                    mutableIntStateOf(0)
                }

                Scaffold(
                    modifier = Modifier.fillMaxSize(),
                    bottomBar = {
                        BottomNavigationBar(
                            selectedItem = selectedItem,
                            onItemSelected = { selectedItem = it }
                        )
                    }
                ) { innerPadding ->
                    Box(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(innerPadding)
                    ) {
                        when (selectedItem) {
                            0 -> CatalogScreen(
                                state = CatalogUiState.Success(sampleArticles),
                                onArticleClick = {},
                                onRetry = {}
                            )
                            1 -> Text("Košarica")
                            2 -> Text("Transakcije")
                            3 -> Text("Račun")
                            else -> Text("Artikli")
                        }
                    }
                }
            }
        }
    }
}