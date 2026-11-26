package hr.foi.air.sevenpay.ui.admin.users

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.KeyboardArrowRight
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hr.foi.air.sevenpay.core.model.Role
import hr.foi.air.sevenpay.core.model.User
import hr.foi.air.sevenpay.ui.theme.Border
import hr.foi.air.sevenpay.ui.theme.DeactivatedRed
import hr.foi.air.sevenpay.ui.theme.InStockText
import hr.foi.air.sevenpay.ui.theme.LightRose
import hr.foi.air.sevenpay.ui.theme.Raspberry
import hr.foi.air.sevenpay.ui.theme.SecondaryText
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme
import hr.foi.air.sevenpay.ui.theme.UnavailableBackground
import hr.foi.air.sevenpay.ui.theme.UnavailableText

@Composable
fun UserListScreen(
    state: UserListUiState,
    onRetry: () -> Unit
) {
    when (state) {
        UserListUiState.Loading -> Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            CircularProgressIndicator()
        }

        is UserListUiState.Error -> Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(24.dp),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(text = state.message, textAlign = TextAlign.Center)
            Spacer(modifier = Modifier.height(16.dp))
            Button(onClick = onRetry) {
                Text("Pokušaj ponovno")
            }
        }

        is UserListUiState.Success -> LazyColumn(
            modifier = Modifier.fillMaxSize(),
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            item {
                Column(modifier = Modifier.padding(bottom = 8.dp)) {
                    Text(
                        text = "Korisnici",
                        fontSize = 28.sp,
                        fontWeight = FontWeight.Bold
                    )
                    Text(
                        text = "Ukupno računa: ${state.users.size}",
                        fontSize = 14.sp,
                        color = SecondaryText
                    )
                }
            }
            items(state.users, key = { it.id }) { user ->
                UserRow(user)
            }
        }
    }
}

@Composable
private fun UserRow(user: User) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(
            containerColor = if (user.active) Color.White else UnavailableBackground.copy(alpha = 0.3f)
        ),
        border = BorderStroke(1.dp, Border)
    ) {
        Row(
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            UserAvatar(user)
            Spacer(modifier = Modifier.width(12.dp))
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = "${user.firstName} ${user.lastName}",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = if (user.active) Color.Unspecified else UnavailableText
                )
                Text(
                    text = user.username,
                    fontSize = 13.sp,
                    color = SecondaryText
                )
                Spacer(modifier = Modifier.height(6.dp))
                Row(
                    horizontalArrangement = Arrangement.spacedBy(10.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    RoleChip(user.role)
                    StatusLabel(user.active)
                }
            }
            Icon(
                imageVector = Icons.AutoMirrored.Filled.KeyboardArrowRight,
                contentDescription = null,
                tint = SecondaryText
            )
        }
    }
}

@Composable
private fun RoleChip(role: Role) {
    val isAdmin = role == Role.ADMIN
    Text(
        text = if (isAdmin) "Administrator" else "Trgovac",
        fontSize = 12.sp,
        fontWeight = FontWeight.Medium,
        color = if (isAdmin) Raspberry else SecondaryText,
        modifier = Modifier
            .background(
                color = if (isAdmin) LightRose else UnavailableBackground.copy(alpha = 0.6f),
                shape = RoundedCornerShape(6.dp)
            )
            .padding(horizontal = 8.dp, vertical = 2.dp)
    )
}

@Composable
private fun StatusLabel(active: Boolean) {
    val color = if (active) InStockText else DeactivatedRed
    Row(verticalAlignment = Alignment.CenterVertically) {
        Box(
            modifier = Modifier
                .size(8.dp)
                .background(color = color, shape = CircleShape)
        )
        Spacer(modifier = Modifier.width(4.dp))
        Text(
            text = if (active) "Aktivan" else "Deaktiviran",
            fontSize = 12.sp,
            fontWeight = FontWeight.Medium,
            color = color
        )
    }
}

@Composable
private fun UserAvatar(user: User) {
    val initials = "${user.firstName.firstOrNull() ?: ""}${user.lastName.firstOrNull() ?: ""}"
    Box(
        modifier = Modifier
            .size(44.dp)
            .background(
                color = if (user.active) LightRose else UnavailableBackground,
                shape = CircleShape
            ),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = initials.uppercase(),
            fontSize = 16.sp,
            fontWeight = FontWeight.Bold,
            color = if (user.active) Raspberry else UnavailableText
        )
    }
}

@Preview(showBackground = true)
@Composable
private fun UserListScreenPreview() {
    SevenPayTheme {
        UserListScreen(
            state = UserListUiState.Success(
                listOf(
                    User(1, "Marko", "Markić", "admin", Role.ADMIN, true),
                    User(2, "Ana", "Kovač", "akovac", Role.MERCHANT, true),
                    User(3, "Ivan", "Horvat", "ihorvat", Role.MERCHANT, false)
                )
            ),
            onRetry = {}
        )
    }
}