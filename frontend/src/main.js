import {createApp} from 'vue'
import App from './App.vue'
import './style.css';

// Disable context menu and pinch/keyboard zoom globally
window.addEventListener('contextmenu', (event) => {
  event.preventDefault();
});

// Block Ctrl/Cmd + mousewheel zoom (WebView2/Electron-like)
window.addEventListener(
  'wheel',
  (event) => {
    if (event.ctrlKey || event.metaKey) {
      event.preventDefault();
    }
  },
  { passive: false }
);

// Block keyboard zoom shortcuts
window.addEventListener('keydown', (event) => {
  const isZoomKey = ['+', '-', '=', '0'].includes(event.key);
  if ((event.ctrlKey || event.metaKey) && isZoomKey) {
    event.preventDefault();
  }
});

createApp(App).mount('#app')
