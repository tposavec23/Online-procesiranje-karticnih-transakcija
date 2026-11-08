
package hr.foi.air.sevenpay.navigation
import androidx.compose.ui.tooling.preview.Preview
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.AccountCircle
import androidx.compose.material.icons.outlined.Receipt
import androidx.compose.material.icons.outlined.ShoppingCart
import androidx.compose.material.icons.outlined.ViewModule
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import hr.foi.air.sevenpay.ui.theme.Raspberry
import hr.foi.air.sevenpay.ui.theme.SecondaryText
import hr.foi.air.sevenpay.ui.theme.CreamBackground
import hr.foi.air.sevenpay.ui.theme.LightRose


@Composable
fun BottomNavigationBar(
    selectedItem: Int,
    onItemSelected: (Int) -> Unit
) {
    val items = listOf(
        Triple("Artikli", Icons.Outlined.ViewModule, 0),
        Triple("Košarica", Icons.Outlined.ShoppingCart, 1),
        Triple("Transakcije", Icons.Outlined.Receipt, 2),
        Triple("Račun", Icons.Outlined.AccountCircle, 3)
    )

    NavigationBar(containerColor = CreamBackground) {
        items.forEach { (title, icon, index) ->
            NavigationBarItem(
                selected = selectedItem == index,
                onClick = { onItemSelected(index) },
                icon = {
                    Icon(
                        imageVector = icon,
                        contentDescription = title
                    )
                },
                label = { Text(title) },
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
            selectedItem = 0,
            onItemSelected = {}
        )
    }
}
