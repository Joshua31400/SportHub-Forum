document.addEventListener('DOMContentLoaded', function() {
    const passwordInput = document.getElementById('password');
    const submitButton = document.querySelector('button[type="submit"]');

    const validationMessage = document.createElement('div');
    validationMessage.className = 'password-validation';
    validationMessage.style.color = '#e53e3e';
    validationMessage.style.fontSize = '14px';
    validationMessage.style.marginTop = '5px';
    validationMessage.style.textAlign = 'left';
    validationMessage.textContent = 'Password must be at least 6 characters long';

    passwordInput.parentNode.insertBefore(validationMessage, passwordInput.nextSibling);

    submitButton.disabled = true;
    submitButton.style.opacity = '0.5';
    submitButton.style.cursor = 'not-allowed';

    function checkPasswordLength() {
        const minLength = 6;
        const currentLength = passwordInput.value.length;

        if (currentLength >= minLength) {
            submitButton.disabled = false;
            submitButton.style.opacity = '1';
            submitButton.style.cursor = 'pointer';
            validationMessage.style.color = '#10b981';
            validationMessage.textContent = 'Password is valid !';
        } else {
            submitButton.disabled = true;
            submitButton.style.opacity = '0.5';
            submitButton.style.cursor = 'not-allowed';
            validationMessage.style.color = '#e53e3e';
            validationMessage.textContent = `Passord nedd to contain ${minLength} caractere (${currentLength}/${minLength})`;
        }
    }
    passwordInput.addEventListener('input', checkPasswordLength);

    checkPasswordLength();
});