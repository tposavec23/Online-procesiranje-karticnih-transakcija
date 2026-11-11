package hr.foi.air.sevenpay.ui.login

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.ErrorOutline
import androidx.compose.material.icons.outlined.Lock
import androidx.compose.material.icons.outlined.Person
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import hr.foi.air.sevenpay.ui.theme.CreamBackground
import hr.foi.air.sevenpay.ui.theme.Raspberry
import hr.foi.air.sevenpay.ui.theme.PrimaryText
import hr.foi.air.sevenpay.ui.theme.Border
import hr.foi.air.sevenpay.ui.theme.SevenPayTheme
import hr.foi.air.sevenpay.ui.theme.SecondaryText
import hr.foi.air.sevenpay.ui.theme.LightRose

@Composable
fun LoginScreen(
    state: LoginUiState,
    onUsernameChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onLoginClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    val hasError = state.errorMessage != null

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(CreamBackground)
            .padding(horizontal = 24.dp, vertical = 32.dp)
    ) {
        AppLogo()
        Spacer(Modifier.height(48.dp))

        Text("Prijava", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = PrimaryText)
        Text("Prijavi se svojim korisničkim podacima.", fontSize = 14.sp, color = SecondaryText)
        Spacer(Modifier.height(24.dp))

        state.errorMessage?.let {
            ErrorBanner(it)
            Spacer(Modifier.height(16.dp))
        }

        LoginField(
            label = "Korisničko ime",
            value = state.username,
            onValueChange = onUsernameChange,
            placeholder = "npr. akovac",
            icon = Icons.Outlined.Person,
            isError = hasError,
            enabled = !state.isLoading
        )
        Spacer(Modifier.height(16.dp))

        LoginField(
            label = "Lozinka",
            value = state.password,
            onValueChange = onPasswordChange,
            placeholder = "Lozinka",
            icon = Icons.Outlined.Lock,
            isError = hasError,
            enabled = !state.isLoading,
            isPassword = true
        )
        Spacer(Modifier.height(24.dp))

        Button(
            onClick = onLoginClick,
            enabled = state.canSubmit,
            shape = RoundedCornerShape(12.dp),
            colors = ButtonDefaults.buttonColors(containerColor = Raspberry),
            modifier = Modifier
                .fillMaxWidth()
                .height(52.dp)
        ) {
            if (state.isLoading) {
                CircularProgressIndicator(color = Color.White, strokeWidth = 2.dp, modifier = Modifier.size(20.dp))
            } else {
                Text("Prijavi se", fontSize = 16.sp, fontWeight = FontWeight.Bold)
            }
        }

        Spacer(Modifier.weight(1f))
        Text(
            "Račun ti otvara administrator.",
            fontSize = 13.sp,
            color = SecondaryText,
            modifier = Modifier.align(Alignment.CenterHorizontally)
        )
    }
}

@Composable
private fun AppLogo() {
    Row(verticalAlignment = Alignment.CenterVertically) {
        Box(
            modifier = Modifier
                .size(60.dp)
                .background(Raspberry, RoundedCornerShape(12.dp)),
            contentAlignment = Alignment.Center
        ) {
            Text("Logo", color = Color.White, fontSize = 22.sp, fontWeight = FontWeight.Bold)
        }
        Spacer(Modifier.width(12.dp))
        Text("7pay", fontSize = 26.sp, fontWeight = FontWeight.Bold, color = PrimaryText)
    }
}

@Composable
private fun ErrorBanner(message: String) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .background(LightRose, RoundedCornerShape(12.dp))
            .padding(12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(Icons.Outlined.ErrorOutline, contentDescription = null, tint = MaterialTheme.colorScheme.error)
        Spacer(Modifier.width(8.dp))
        Text(message, fontSize = 14.sp, color = MaterialTheme.colorScheme.error)
    }
}

@Composable
private fun LoginField(
    label: String,
    value: String,
    onValueChange: (String) -> Unit,
    placeholder: String,
    icon: ImageVector,
    isError: Boolean,
    enabled: Boolean,
    isPassword: Boolean = false
) {
    Column {
        Text(label, fontSize = 14.sp, fontWeight = FontWeight.SemiBold, color = PrimaryText)
        Spacer(Modifier.height(6.dp))
        OutlinedTextField(
            value = value,
            onValueChange = onValueChange,
            placeholder = { Text(placeholder, color = SecondaryText) },
            leadingIcon = { Icon(icon, contentDescription = null, tint = SecondaryText) },
            isError = isError,
            enabled = enabled,
            singleLine = true,
            visualTransformation = if (isPassword) PasswordVisualTransformation() else VisualTransformation.None,
            keyboardOptions = KeyboardOptions(
                keyboardType = if (isPassword) KeyboardType.Password else KeyboardType.Text
            ),
            shape = RoundedCornerShape(12.dp),
            colors = OutlinedTextFieldDefaults.colors(
                focusedContainerColor = Color.White,
                unfocusedContainerColor = Color.White,
                disabledContainerColor = Color.White,
                unfocusedBorderColor = Border,
                focusedBorderColor = Raspberry
            ),
            modifier = Modifier.fillMaxWidth()
        )
    }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun LoginScreenPreview() {
    SevenPayTheme {
        LoginScreen(LoginUiState(), {}, {}, {})
    }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun LoginScreenErrorPreview() {
    SevenPayTheme {
        LoginScreen(
            LoginUiState(
                username = "akovac",
                password = "pogresna",
                errorMessage = LoginErrorMessages.INVALID_CREDENTIALS
            ),
            {}, {}, {}
        )
    }
}

@Preview(showBackground = true, heightDp = 800)
@Composable
private fun LoginScreenLoadingPreview() {
    SevenPayTheme {
        LoginScreen(LoginUiState(username = "akovac", password = "trgovac123", isLoading = true), {}, {}, {})
    }
}