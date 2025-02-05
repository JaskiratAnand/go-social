import { BrowserRouter, Routes, Route } from "react-router"
import ConfirmationPage from "./ConfirmationPage"

export const API_URL = import.meta.env.VITE_API_URL as string || 'http://localhost:8080/v1';
function App() {

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<div> Home </div>} />
        <Route path="/confirm/:token" element={<ConfirmationPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
