import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import hr.foi.air.sevenpay.ui.theme.CreamBackground
import hr.foi.air.sevenpay.ui.theme.LightRose
import hr.foi.air.sevenpay.ui.theme.Raspberry
import hr.foi.air.sevenpay.ui.theme.SecondaryText

@Composable
fun AdminBottomNavigationBar(
    selectedTab: AdminTab,
    onTabSelected: (AdminTab) -> Unit
) {
    NavigationBar(containerColor = CreamBackground) {
        AdminTab.entries.forEach { tab ->
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