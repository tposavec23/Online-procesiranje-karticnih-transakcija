package hr.foi.air.sevenpay.navigation

import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.AccountCircle
import androidx.compose.material.icons.outlined.Receipt
import androidx.compose.material.icons.outlined.ShoppingCart
import androidx.compose.material.icons.outlined.ViewModule
import androidx.compose.ui.graphics.vector.ImageVector

// Tabs of the merchant's bottom navigation bar, in the order they are shown.
enum class MerchantTab(val title: String, val icon: ImageVector) {
    ARTICLES("Artikli", Icons.Outlined.ViewModule),
    CART("Košarica", Icons.Outlined.ShoppingCart),
    TRANSACTIONS("Transakcije", Icons.Outlined.Receipt),
    ACCOUNT("Račun", Icons.Outlined.AccountCircle)
}
