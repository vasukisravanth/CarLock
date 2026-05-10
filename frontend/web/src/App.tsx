import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import LockControl from './components/LockControl';

const App: React.FC = () => {
  return (
    <Router>
      <div>
        <h1>Car Lock System</h1>
        <Switch>
          <Route path="/" exact component={LockControl} />
        </Switch>
      </div>
    </Router>
  );
};

export default App;