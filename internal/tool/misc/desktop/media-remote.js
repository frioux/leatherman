var enabled = true

function over(x) {
   if (!enabled) { return }
   let formData = new FormData();

   formData.append('player', x.innerText);

   fetch('/play', {
      method: 'POST',
      body: formData
   })
      .then(result => {
         console.log('Success:', result);
      })
      .catch(error => {
         console.error('Error:', error);
      });

}

function out(x) {
   if (!enabled) { return }
   let formData = new FormData();

   formData.append('player', x.innerText);

   fetch('/pause', {
      method: 'POST',
      body: formData,
   })
      .then(result => {
         console.log('Success:', result);
      })
      .catch(error => {
         console.error('Error:', error);
      });
}

function selectOutput(x) {
   if (!enabled) { return }
   enabled = false
   let formData = new FormData();

   formData.append('player', x.innerText);

   var myHeaders = new Headers();
   myHeaders.append('Accept', 'text/html');

   fetch('/select-player', {
      method: 'POST',
      body: formData,
      headers: { 'Accept': 'application/json' }, // work around fraidycat workaround
   })
      .then(result => {
         console.log('Success:', result);
      })
      .catch(error => {
         console.error('Error:', error);
      });
}

