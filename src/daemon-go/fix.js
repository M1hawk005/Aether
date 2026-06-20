const fs = require('fs');
['storage.go', 'api.go'].forEach(f => {
  let c = fs.readFileSync(f, 'utf8');
  c = c.replace(/\\\`/g, '\`');
  fs.writeFileSync(f, c);
});
