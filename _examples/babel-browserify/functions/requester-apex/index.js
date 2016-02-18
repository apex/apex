
import axios from 'axios'
import λ from 'apex.js'
import 'babel-polyfill'

export default λ(e => {
  console.log('fetching %d urls', e.urls.length)
  return Promise.all(e.urls.map(async function(url){
    console.log('fetching %s', url)
    return {
      status: (await axios.get(url)).status,
      url
    }
  }))
})
