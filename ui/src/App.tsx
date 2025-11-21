import { Routes, Route } from 'react-router-dom'
import { Toaster } from 'sonner'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Clusters from './pages/Clusters'
import CreateCluster from './pages/CreateCluster'
import Services from './pages/Services'
import Monitoring from './pages/Monitoring'
import Routing from './pages/Routing'
import Settings from './pages/Settings'

function App() {
  return (
    <>
      <Routes>
        {/* Full-page routes without layout */}
        <Route path="/clusters/create" element={<CreateCluster />} />
        
        {/* Routes with layout */}
        <Route path="*" element={
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
        } />
      </Routes>
      <Toaster position="top-right" richColors />
    </>
  )
}

export default App

