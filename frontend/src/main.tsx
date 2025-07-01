import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { ToastContainer } from "react-fox-toast"
import './index.css'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    {/* @ts-ignore */}
    <ToastContainer toastTypeTheming={{
    success: {
      style: {
        backgroundColor: '#3b3b3b',
        color: '#e6e6e6',
      },
      className: '',
    },
    error: {
      style: {
        backgroundColor: '#3b3b3b',
        color: '#e6e6e6',
      },
      className: ''
    }}}/> 
    <App />
  </StrictMode>,
)
