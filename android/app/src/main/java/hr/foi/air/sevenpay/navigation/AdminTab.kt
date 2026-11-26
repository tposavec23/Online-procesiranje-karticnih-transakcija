import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.AccountCircle
import androidx.compose.material.icons.outlined.Category
import androidx.compose.material.icons.outlined.Group
import androidx.compose.material.icons.outlined.ViewModule
import androidx.compose.ui.graphics.vector.ImageVector

enum class AdminTab(val title: String, val icon: ImageVector) {
    ARTICLES("Artikli", Icons.Outlined.ViewModule),
    CATEGORIES("Kategorije", Icons.Outlined.Category),
    USERS("Korisnici", Icons.Outlined.Group),
    ACCOUNT("Račun", Icons.Outlined.AccountCircle)
}