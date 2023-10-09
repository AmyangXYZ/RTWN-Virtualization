import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import Vuesax from 'vuesax3'

import './assets/main.css'
import 'vuesax3/dist/vuesax.css'
import 'material-icons/iconfont/material-icons.css'

import { fetchNetwork } from './hooks/useNetwork.js'
await fetchNetwork()

const app = createApp(App)

app.use(router)
app.use(Vuesax)
app.mount('#app')
