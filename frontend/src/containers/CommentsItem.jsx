import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';
import emotes from '../utils/emotes.json';

import Markdown from '../components/Markdown';
import Textarea from '../components/Textarea';

@connect(createStructuredSelector({
  myId: selectors.myId,
  amAdmin: selectors.amAdmin,
  users: selectors.users,
}), actions)
export default class CommentsItem extends Component {
  static propTypes = {
    comment: ImmutablePropTypes.map,
    userId: PropTypes.string,
    glue: PropTypes.bool,

    myId: PropTypes.string,
    amAdmin: PropTypes.bool,
    users: ImmutablePropTypes.map,

    updateComment: PropTypes.func.isRequired,
    patchComment: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  state = {
    editing: false,
    peeking: false,
  };

  componentWillReceiveProps({ comment }) {
    if (comment === this.props.comment) return;
    if (!comment || comment.get('isDeleted')) this.setState({ editing: false });
  }

  handleEdit = () => this.setState({ editing: true });

  handleEditSave = () => this.props.updateComment(
    this.props.comment.get('id'),
    { body: this.r.input.state.value },
    () => {
      this.setState({ editing: false });
      this.props.messagePush('updated a comment');
    },
  );

  handleEditDiscard = () => {
    if (!window.confirm('Discard your edits?')) return;
    this.setState({ editing: false });
  };

  handleDelete = () => {
    if (!this.props.comment) return;
    const { isDeleted, id } = this.props.comment.toJS();
    if (!isDeleted && !window.confirm('Delete this comment?')) return;
    this.props.patchComment(id, { action: isDeleted ? 'undelete' : 'delete' });
  };

  handlePeek = () => this.setState({ peeking: !this.state.peeking });

  render() {
    const { comment, userId, glue, myId, amAdmin, users } = this.props;
    if (!comment) {
      return (
        <div className={this.cni()}>
          <div className={this.cn({ d: 'body', m: 'noComments' })}>
            No comments
          </div>
        </div>
      );
    }

    const x = comment.toJS();
    const user = users.get(userId);
    const [isAdmin, name, gravatar] = user ? [
      user.get('isAdmin'),
      user.get('nickname'),
      user.get('gravatar'),
    ] : [null, null, null];
    const createdAt = moment(x.createdAt);
    const updatedAt = x.updatedAt && moment(x.updatedAt);

    let { body } = x;
    if (x.isDeleted && (!amAdmin || !this.state.peeking)) {
      body = 'DELETED ¯\\\\\\_(ツ)_/¯';
    } else {
      for (const [emote, { image_id: id }] of Object.entries(emotes.emotes)) {
        body = body.replace(
          new RegExp(emote, 'g'),
          `![${emote}](${emotes.template.small.replace('{image_id}', id)})`,
        );
      }
    }

    return (
      <div className={this.cni({
        m: [
          glue && 'glue',
          isAdmin && 'admin',
        ],
        s: [
          x.isDeleted && 'deleted',
          this.state.editing && 'editing',
        ],
      })}>
        {!glue && <Link
          className={this.cn({ d: 'author' })}
          to={`/users/${x.createdBy}`}
        >
          <div className={this.cn({
            d: 'authorName',
            m: isAdmin && 'withRole',
          })}>
            {name}{isAdmin && <div className={this.cn({ d: 'authorRole' })}>
              admin
            </div>}
          </div>
          <img
            className={this.cn({ d: 'authorAvatar' })}
            src={`//gravatar.com/avatar/${gravatar}?s=46&d=identicon`}
          />
        </Link>}
        <div className={this.cn({
          d: 'meta',
          m: updatedAt && 'withEdited',
        })}>
          <span
            className={this.cn({
              d: 'metaCreatedAt',
              m: updatedAt && 'withEdited',
            })}
            title={createdAt.format('lll')}
          >
            {createdAt.fromNow()}
          </span>
          {!this.state.editing && myId === x.createdBy && !x.isDeleted &&
            <button
              className={this.cn({ d: 'metaAction' })}
              title="edit"
              onClick={this.handleEdit}
            >✎</button>
          }
          {this.state.editing && <button
            className={this.cn({ d: 'metaAction' })}
            title="save"
            onClick={this.handleEditSave}
          >✔</button>}
          {this.state.editing && <button
            className={this.cn({ d: 'metaAction' })}
            title="discard"
            onClick={this.handleEditDiscard}
          >✘</button>}
          {((myId === x.createdBy && !x.isDeleted) || amAdmin) && <button
            className={this.cn({ d: 'metaAction' })}
            title={x.isDeleted ? 'undelete' : 'delete'}
            onClick={this.handleDelete}
          >{x.isDeleted ? '♻' : '×'}</button>}
          {x.isDeleted && amAdmin && <button
            className={this.cn({ d: 'metaAction' })}
            title={!this.state.peeking ? 'see original' : 'hide original'}
            onClick={this.handlePeek}
          >👁</button>}
          {updatedAt && <div
            className={this.cn({ d: 'metaEdited' })}
            title={updatedAt.format('lll')}
          >
            {x.isDeleted ? 'deleted' : 'last edited'} <b>
              {updatedAt.fromNow()}
            </b>
          </div>}
        </div>
        {!this.state.editing && <Markdown
          className={this.cn({
            d: 'body',
            m: x.isDeleted && `deleted${!this.state.peeking ? 'Snip' : ''}`,
          })}
          source={body}
        />}
        {this.state.editing && <Textarea
          ref={this.rcb('input')}
          className={this.cn({ d: 'editor' })}
          value={x.body}
          label="Comment editor"
        />}
      </div>
    );
  }
}
