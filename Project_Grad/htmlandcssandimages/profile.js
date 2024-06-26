document.addEventListener('DOMContentLoaded', () => {
    fetch('/api/getUserData')
        .then(response => response.json())
        .then(data => {
            document.getElementById('name').innerText = data.name;
            document.getElementById('email').value = data.email;
            document.getElementById('password').value = data.password;
            document.getElementById('phone').value = data.phone;
        })
        .catch(error => console.error('Error fetching user data:', error));
});