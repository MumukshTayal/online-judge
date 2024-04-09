import './App.css' 
import { LoginPage } from "@/components/login-page.tsx"
import { AddTestCasePage } from "@/components/add-test-case-page.tsx"
import { SubmissionsPage } from "@/components/submissions_page.tsx"
import { GoogleOAuthProvider } from '@react-oauth/google'; 


function App() {
  const key = "461833529604-i7difeupuv73o6v8j1jc0s75lr44ms6u.apps.googleusercontent.com"
  return (
    <GoogleOAuthProvider clientId={key}>
      <>
        <LoginPage />
      </>
    </GoogleOAuthProvider>
  )
}

export default App
