// script.js

document.addEventListener('DOMContentLoaded', () => {
    const toggleButton = document.getElementById('dark-mode-toggle');
    const darkModeIcon = document.getElementById('dark-mode-icon');
    
    // Load the saved theme from local storage, if available
    if (localStorage.getItem('theme') === 'dark') {
        document.body.classList.add('dark-mode');
        darkModeIcon.src = '/Users/khaled/Desktop/Project Grad/images/sun.png';
    } else {
        darkModeIcon.src = '/Users/khaled/Desktop/Project Grad/images/moon.png';
    }

    toggleButton.addEventListener('click', () => {
        document.body.classList.toggle('dark-mode');
        
        if (document.body.classList.contains('dark-mode')) {
            localStorage.setItem('theme', 'dark');
            darkModeIcon.src = '/Users/khaled/Desktop/Project Grad/images/sun.png';
        } else {
            localStorage.removeItem('theme');
            darkModeIcon.src = '/Users/khaled/Desktop/Project Grad/images/moon.png';
        }
    });
});
