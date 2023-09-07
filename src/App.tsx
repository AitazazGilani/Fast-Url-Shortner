import { Button } from "./components/button"
import { Input } from "./components/input"
import { Card } from "./components/Card"


function App() {
  return (
    <>
    
  <div className="">
  <h1 className="right-20">Lightning Fast Url Shortner ⚡</h1>
  <p className="leading-7 [&:not(:first-child)]:mt-6 text-lg font-semibold">
      Tired of long amazon links? try shortning them!
    </p>
  <Input className="" placeholder="Link goes here" />
  <Button> Next</Button>
  </div>
    </>
  )
}

export default App
