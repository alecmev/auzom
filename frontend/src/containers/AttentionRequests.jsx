import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

const attentionRequestsSelector = createSelector(
  selectors.attentionRequests,
  (_, props) => props.matchId,
  (x, matchId) =>
    matchId && x.toList().filter(y =>
      y.get('target') === 'match' &&
      y.get('targetId') === matchId,
    ),
);

const pendingSelector = createSelector(
  attentionRequestsSelector,
  x => x && x.filter(y => y.get('claimedAt') === null),
);

const claimedSelector = createSelector(
  attentionRequestsSelector,
  x => x && x.filter(y =>
    y.get('claimedAt') !== null &&
    y.get('resolvedAt') === null,
  ),
);

const resolvedSelector = createSelector(
  attentionRequestsSelector,
  x => x && x.filter(y => y.get('resolvedAt') !== null),
);

@connect(createStructuredSelector({
  attentionRequests: attentionRequestsSelector,
  pending: pendingSelector,
  claimed: claimedSelector,
  resolved: resolvedSelector,
  users: selectors.users,
}), actions)
export default class AttentionRequests extends Component {
  static propTypes = {
    matchId: PropTypes.string,

    attentionRequests: ImmutablePropTypes.list,
    pending: ImmutablePropTypes.list,
    claimed: ImmutablePropTypes.list,
    resolved: ImmutablePropTypes.list,
    users: ImmutablePropTypes.map,

    loadAttentionRequests: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ matchId, attentionRequests }, prevProps) {
    if (!matchId) return;
    if (
      !prevProps || matchId !== prevProps.matchId ||
      (!attentionRequests && prevProps.attentionRequests)
    ) {
      this.props.loadAttentionRequests({
        target: 'match',
        targetId: matchId,
      }, undefined, ['createdBy', 'claimedBy']);
    }
  }

  render() {
    const { attentionRequests, pending, claimed, resolved, users } = this.props;
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          {attentionRequests && !attentionRequests.size && 'No requests'}
          {pending && pending.size && <div>Pending:
            <ul>{pending.toJS().map(x => (
              <li key={x.id}>
                [{users.getIn([x.createdBy, 'nickname'])}] {x.message}
              </li>
            ))}</ul>
          </div>}
          {claimed && !!claimed.size && <div>In progress:
            <ul>{claimed.toJS().map(x => (
              <li key={x.id}>
                [{users.getIn([x.createdBy, 'nickname'])}] {x.message}
              </li>
            ))}</ul>
          </div>}
          {resolved && !!resolved.size && <div>Resolved:
            <ul>{resolved.toJS().map(x => (
              <li key={x.id}>
                {x.isDiscarded ? '[DISCARDED] ' : ''}
                [{users.getIn([x.createdBy, 'nickname'])}] {x.message}
              </li>
            ))}</ul>
          </div>}
        </div>
      </div>
    );
  }
}
