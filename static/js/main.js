// Main JavaScript for landing and writing pages

// Hamburger menu toggle
document.addEventListener('DOMContentLoaded', function() {
    const hamburgerBtn = document.getElementById('hamburgerBtn');
    const dropdownMenu = document.getElementById('dropdownMenu');
    
    if (hamburgerBtn && dropdownMenu) {
        hamburgerBtn.addEventListener('click', function(e) {
            e.stopPropagation();
            dropdownMenu.classList.toggle('show');
        });
        
        // Close dropdown when clicking outside
        document.addEventListener('click', function(e) {
            if (!hamburgerBtn.contains(e.target) && !dropdownMenu.contains(e.target)) {
                dropdownMenu.classList.remove('show');
            }
        });
    }
    
    // Account dropdown toggle
    const accountBtn = document.getElementById('accountBtn');
    const accountDropdown = document.getElementById('accountDropdown');
    
    if (accountBtn && accountDropdown) {
        accountBtn.addEventListener('click', function(e) {
            e.stopPropagation();
            accountDropdown.classList.toggle('show');
        });
        
        // Close dropdown when clicking outside
        document.addEventListener('click', function(e) {
            if (!accountBtn.contains(e.target) && !accountDropdown.contains(e.target)) {
                accountDropdown.classList.remove('show');
            }
        });
    }
});

