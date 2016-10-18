import { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import * as actions from '../actions';

import New from './New';

@connect(undefined, actions)
export default class Signup extends Component {
  static propTypes = {
    createUser: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  handleSuccess = () => {
    browserHistory.push('/');
    this.props.messagePush('signed up, check your email');
  };

  render() {
    return (
      <New
        fields={{
          email: { type: 'email', name: 'Email' },
          password: { type: 'password', name: 'Password' },
        }}
        callForAction="sign up"
        link={{ to: '/login', name: 'Already have an account?' }}
        createResource={this.props.createUser}
        onSuccess={this.handleSuccess}
        noMessage
      />
    );
  }
}
