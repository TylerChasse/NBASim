console.log("JavaScript file loaded successfully!");

dropdowns = document.querySelectorAll('.dropdown');
dropdowns.forEach(dropdown => {
    const select = dropdown.querySelector('.select');
    const caret = dropdown.querySelector('.caret');
    const menu = dropdown.querySelector('.menu');
    const options = dropdown.querySelectorAll('.menu li');
    const selected = dropdown.querySelector('.selected');

    select.addEventListener('click', () => {
        select.classList.toggle('select-clicked');
        caret.classList.toggle('caret-rotate');
        menu.classList.toggle('menu-open');
    });

    options.forEach(option => {
        option.addEventListener('click', () => {
            selected.innerText = option.innerText;
            select.classList.remove('select-clicked');
            caret.classList.remove('caret-rotate');
            menu.classList.remove('menu-open');

            options.forEach(option => {
                option.classList.remove('active');
            });
            option.classList.add('active');
        });
    });
});

document.addEventListener("DOMContentLoaded", function () {
    var form = document.getElementById("go");

    if (form) {
        form.addEventListener("submit", function (event) {
            event.preventDefault(); // Prevent the form from submitting before modifying the action
            editFormAction();
            form.submit(); // Now submit with the correct action
        });
    }
});

function editFormAction() {
    var form = document.getElementById('go');

    var activeElements = document.getElementsByClassName('active');

    if (!form) {
        console.error("Form with ID 'go' not found.");
        return;
    }

    if (activeElements.length < 2) {
        console.error("Not enough active elements selected.");
        return;
    }

    var team1Abb = document.getElementsByClassName('active')[0].innerText;
    var team2Abb = document.getElementsByClassName('active')[1].innerText;
    var start = form.getAttribute("data-start");
    form.action = start.concat(team1Abb, "/", team2Abb);
};