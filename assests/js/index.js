window.addEventListener('load', () => {
    const url    = document.querySelector('#url');
    const group  = document.querySelector('#group');
    const addURL = document.querySelector('#add-url');
    const urls   = document.querySelector('#add-url');

    updateURLs();

    addURL.addEventListener('submit', (ev) => {
        ev.preventDefault();
        if (!url.value) return false;
        const fetchParams = {
            method: "post",
            headers: {
                "Content-type": "application/json"
            },
            body: JSON.stringify({ url: url.value, group: group.value })
        };
        fetch('/add', fetchParams)
            .then(res => updateURLs())
            .catch(err => console.error(err));
        url.value = group.value = '';
    });
});

function updateURLs() {
    urls.innerHTML = "";
    fetch('/urls')
        .then(res => res.json())
        .then(data => {
            let html = "";
            for (const group in data) {
                const block = document.createElement('table');
                block.classList.add('u-full-width');
                block.innerHTML += `<thead><th style="text-align: center;">${group}</th></thead>`;
                for (const url of data[group]) {
                    block.innerHTML += `<tr><td><a href="${url}" target="blank">${url}</a></td></tr>`;
                }
                urls.appendChild(block);
            }
        });
}

