import Component from '../utils/Component';

export default class Loading extends Component {
  render() {
    return (
      <div className={this.cni()}>
        <img
          className={this.cn({ d: 'logo' })}
          src={require('../assets/logo.svg')}
        />
      </div>
    );
  }
}
