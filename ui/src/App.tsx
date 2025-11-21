import { Routes, Route } from 'react-router-dom'
import { Toaster } from 'sonner'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Clusters from './pages/Clusters'
import Services from './pages/Services'
import Monitoring from './pages/Monitoring'
import Routing from './pages/Routing'
import Settings from './pages/Settings'

function App() {
  return (
    <>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/clusters" element={<Clusters />} />
          <Route path="/services" element={<Services />} />
          <Route path="/monitoring" element={<Monitoring />} />
          <Route path="/routing" element={<Routing />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Layout>
      <Toaster position="top-right" richColors />
    </>
  )
}

export default App

