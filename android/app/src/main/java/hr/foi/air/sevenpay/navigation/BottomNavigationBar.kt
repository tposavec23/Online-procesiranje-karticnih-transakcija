package hr.foi.air.sevenpay.navigation

import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.tooling.preview.Preview
import hr.foi.air.sevenpay.ui.theme.CreamBackground
import hr.foi.air.sevenpay.ui.theme.LightRose
import hr.foi.air.sevenpay.ui.theme.Raspberry
import hr.foi.air.sevenpay.ui.theme.SecondaryText
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme

@Composable
fun BottomNavigationBar(
    selectedTab: MerchantTab,
    onTabSelected: (MerchantTab) -> Unit
) {
    NavigationBar(containerColor = CreamBackground) {
        MerchantTab.entries.forEach { tab ->
            NavigationBarItem(
                selected = tab == selectedTab,
                onClick = { onTabSelected(tab) },
                icon = {
                    Icon(
                        imageVector = tab.icon,
                        contentDescription = tab.title
                    )
                },
                label = { Text(tab.title) },
                colors = NavigationBarItemDefaults.colors(
                    selectedIconColor = Raspberry,
                    selectedTextColor = Raspberry,
                    unselectedIconColor = SecondaryText,
                    unselectedTextColor = SecondaryText,
                    indicatorColor = LightRose
                )
            )
        }
    }
}

@Preview(showBackground = true)
@Composable
fun BottomNavigationBarPreview() {
    SevenPayTheme {
        BottomNavigationBar(
            selectedTab = MerchantTab.ARTICLES,
            onTabSelected = {}
        )
    }
}
