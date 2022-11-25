import Route from '@ember/routing/route';
//import { service } from '@ember/service';
import { getOwner } from '@ember/application';

export default class IndexRoute extends Route {
  //@service store;
  get store() {
    return getOwner(this).lookup('service:store');
  }

  async model() {
    return this.store.findAll('rental');
  }
}
