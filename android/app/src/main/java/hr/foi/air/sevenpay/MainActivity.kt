
package hr.foi.air.sevenpay

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import hr.foi.air.sevenpay.core.model.Role
import hr.foi.air.sevenpay.navigation.BottomNavigationBar
import hr.foi.air.sevenpay.ui.catalog.CatalogScreen
import hr.foi.air.sevenpay.ui.catalog.CatalogViewModel
import hr.foi.air.sevenpay.ui.login.LoginScreen
import hr.foi.air.sevenpay.ui.login.LoginViewModel
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContent {
            SevenPayTheme {
                val loginViewModel: LoginViewModel = viewModel()
                val loginState by loginViewModel.state.collectAsStateWithLifecycle()

                when (loginState.loggedInUser?.role) {
                    null -> LoginScreen(
                        state = loginState,
                        onUsernameChange = loginViewModel::onUsernameChange,
                        onPasswordChange = loginViewModel::onPasswordChange,
                        onLoginClick = loginViewModel::login
                    )
                    Role.MERCHANT -> MerchantHome()
                    Role.ADMIN -> AdminHome()
                }
            }
        }
    }
}

@Composable
private fun MerchantHome() {
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
                0 -> {
                    val catalogViewModel: CatalogViewModel = viewModel()
                    val state by catalogViewModel.state.collectAsStateWithLifecycle()

                    CatalogScreen(
                        state = state,
                        onArticleClick = { article ->
                            // Article details will be added later.
                        },
                        onRetry = {
                            catalogViewModel.loadArticles()
                        },
                        onLoadMore = {
                            catalogViewModel.loadMore()
                        }
                    )
                }
                1 -> Text("Košarica")
                2 -> Text("Transakcije")
                3 -> Text("Račun")
                else -> Text("Artikli")
            }
        }
    }
}

// Placeholder until the administrator screens are defined.
@Composable
private fun AdminHome() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text("Administrator")
    }
}
