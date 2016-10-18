import { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

import * as actions from '../actions';

import New from './New';

@connect(undefined, actions)
export default class AttentionRequestNew extends Component {
  static propTypes = {
    matchId: PropTypes.string,
    linkPrefix: PropTypes.string,

    createAttentionRequest: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleSuccess = () => {
    this.props.linkPrefix && this.context.router.push(
      `${this.props.linkPrefix}/attention-requests`,
    );
    this.props.messagePush('submitted an admin attention request');
  };

  render() {
    if (!this.props.matchId) return null;
    const target = 'match';
    const targetId = this.props.matchId;
    return (
      <New
        fields={{
          message: { name: 'Message' },
          target: { type: 'hidden', value: target },
          targetId: { type: 'hidden', value: targetId },
        }}
        callForAction="submit attention request"
        createResource={this.props.createAttentionRequest}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
