const express = require('express');
const fs      = require('fs');
const router  = express.Router();


/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('index', { title: 'Express' });
});

router.post('/add', (req, res) => {
  const entry = req.body;

  fs.readFile('./urls.json', 'utf8', (err, data) => {
    if (err) throw err;

    const modifiedData = JSON.parse(data);
    if(!modifiedData[entry.group]) modifiedData[entry.group] = [];
    modifiedData[entry.group].push(entry.url);

    fs.writeFile('./urls.json', JSON.stringify(modifiedData), (err) => {
      if (err) throw err;

      res.send('Done.');
    });
  });
});

router.get('/urls', (req, res) => {
  fs.readFile('./urls.json', 'utf8', (err, data) => {
    res.json(JSON.parse(data)); // I know this is dumb, but it adds Content-type automatically 😅
  });
});

module.exports = router;
