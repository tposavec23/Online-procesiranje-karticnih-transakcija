package hr.foi.air.sevenpay

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.BackHandler
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import hr.foi.air.sevenpay.core.model.Role
import hr.foi.air.sevenpay.navigation.BottomNavigationBar
import hr.foi.air.sevenpay.navigation.MerchantTab
import hr.foi.air.sevenpay.ui.catalogue.detail.ArticleDetailRoute
import hr.foi.air.sevenpay.ui.catalogue.list.CatalogueScreen
import hr.foi.air.sevenpay.ui.catalogue.list.CatalogueViewModel
import hr.foi.air.sevenpay.ui.login.LoginScreen
import hr.foi.air.sevenpay.ui.login.LoginViewModel
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContent {
            SevenPayTheme {
                val loginViewModel: LoginViewModel = viewModel(factory = LoginViewModel.Factory)
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
    var selectedTab by rememberSaveable {
        mutableStateOf(MerchantTab.ARTICLES)
    }
    var selectedArticleId by rememberSaveable {
        mutableStateOf<Int?>(null)
    }

    Scaffold(
        modifier = Modifier.fillMaxSize(),
        bottomBar = {
            BottomNavigationBar(
                selectedTab = selectedTab,
                onTabSelected = {
                    selectedTab = it
                    selectedArticleId = null
                }
            )
        }
    ) { innerPadding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(innerPadding)
        ) {
            when (selectedTab) {
                MerchantTab.ARTICLES -> {
                    val articleId = selectedArticleId

                    if (articleId != null) {
                        BackHandler { selectedArticleId = null }
                        ArticleDetailRoute(
                            articleId = articleId,
                            onBack = { selectedArticleId = null }
                        )
                    } else {
                        val catalogueViewModel: CatalogueViewModel = viewModel(factory = CatalogueViewModel.Factory)
                        val state by catalogueViewModel.state.collectAsStateWithLifecycle()

                        CatalogueScreen(
                            state = state,
                            onArticleClick = { article ->
                                selectedArticleId = article.id
                            },
                            onRetry = {
                                catalogueViewModel.loadArticles()
                            },
                            onLoadMore = {
                                catalogueViewModel.loadMore()
                            }
                        )
                    }
                }
                MerchantTab.CART, MerchantTab.TRANSACTIONS, MerchantTab.ACCOUNT -> Text(selectedTab.title)
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