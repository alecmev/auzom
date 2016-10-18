import Immutable from 'immutable';
import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

import CommentsItem from './CommentsItem';

import Button from '../components/Button';
import Textarea from '../components/Textarea';

const commentsSelector = createSelector(
  selectors.comments,
  (_, props) => props.target,
  (_, props) => props.targetId,
  (x, target, targetId) =>
    target && targetId && x.toList().filter(y =>
      y.get('target') === target &&
      y.get('targetId') === targetId,
    ).sortBy(y => y.get('createdAt')),
);

const userIdsSelector = createSelector(
  commentsSelector,
  x => x && x.reduce((r, y) =>
    r.union([y.get('createdBy'), y.get('updatedBy')]), new Immutable.Set(),
  ),
);

@connect(createStructuredSelector({
  myId: selectors.myId,
  comments: commentsSelector,
  userIds: userIdsSelector,
}), actions)
export default class Comments extends Component {
  static propTypes = {
    target: PropTypes.string,
    targetId: PropTypes.string,

    myId: PropTypes.string,
    comments: ImmutablePropTypes.list,
    userIds: ImmutablePropTypes.set,

    loadComments: PropTypes.func.isRequired,
    loadUser: PropTypes.func.isRequired,
    patchComment: PropTypes.func.isRequired,
    createComment: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ target, targetId, comments, userIds }, prevProps) {
    if (!target || !targetId) return;
    if (
      !prevProps ||
      target !== prevProps.target ||
      targetId !== prevProps.targetId ||
      (!comments && prevProps.comments)
    ) {
      this.props.loadComments({
        target,
        targetId,
      });
    }

    if (!userIds) return;
    if (!prevProps || userIds !== prevProps.userIds) {
      userIds.forEach(x => x && this.props.loadUser(x));
    }
  }

  refresh = () => this.load(this.props);

  submit = () => {
    const { target, targetId, createComment, messagePush } = this.props;
    if (!target || !targetId) return; // TODO: show message or something
    createComment({
      target,
      targetId,
      body: this.r.input.state.value,
    }, () => {
      this.load(this.props);
      messagePush('posted a comment');
      this.r.input.reset();
    });
  };

  render() {
    const { myId, comments } = this.props;
    let prevCreatedBy;
    let prevCreatedAt;
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ d: 'inner' })}>
          {comments && !comments.size && <CommentsItem />}
          {comments && comments.map((x) => {
            const createdBy = x.get('createdBy');
            const createdAt = moment(x.get('createdAt'));
            let glue = false;
            if (createdBy !== prevCreatedBy) {
              prevCreatedBy = createdBy;
            } else if (
              prevCreatedAt &&
              moment.duration(createdAt.diff(prevCreatedAt)).asMinutes() < 5
            ) {
              glue = true;
            }

            prevCreatedAt = createdAt;
            return (
              <CommentsItem
                key={x.get('id')}
                comment={x}
                userId={createdBy}
                glue={glue}
              />
            );
          }).toJS()}
          {myId && <div className={this.cn({ d: 'newCommentOuter' })}>
            <button
              className={this.cn({ d: 'refresh' })}
              onClick={this.refresh}
            >
              <span className={this.cn({ d: 'refreshArrow' })}>⟳</span>
              <span className={this.cn({ d: 'refreshAction' })}>refresh</span>
            </button>
            <Textarea
              ref={this.rcb('input')}
              className={this.cn({ d: 'newComment' })}
              label="New comment"
            />
          </div>}
          {myId && <div className={this.cn({ d: 'submitOuter' })}>
            <Button
              type="important"
              text="post comment"
              onClick={this.submit}
            />
          </div>}
        </div>
      </div>
    );
  }
}
