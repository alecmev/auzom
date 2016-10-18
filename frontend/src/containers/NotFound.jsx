import Component from '../utils/Component';

export default class NotFound extends Component {
  render() {
    return (
      <div className={this.cni()}>
        404
        <div className={this.cn({ d: 'sub' })}>Not Found ;(</div>
      </div>
    );
  }
}
