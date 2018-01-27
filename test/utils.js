import imgMarioRunning from "./sprites/marioRunning.png";
import imgMarioIdle from "./sprites/marioIdle.png";
import imgMarioJumping from "./sprites/marioJumping.png";
import imgGroundBlock from "./sprites/groundBlock.png";
import imgCoin from "./sprites/coin.png";
import imgGoomba from "./sprites/goomba.png";
import imgQuestionBlock from "./sprites/questionBlock.png";

function unorderedList() {
  return {
    items: [],
    freeSpaces: [],
    add(value) {
      let index = this.items.length;
      if (this.freeSpaces.length) {
        index = this.freeSpaces.pop();
      }
      this.items[index] = value;
      return index;
    },
    remove(index) {
      this.freeSpaces.push(index);
      this.items[index] = unorderedList.REMOVED_ITEM;
    }
  };
}
unorderedList.REMOVED_ITEM = Symbol("REMOVED_ITEM");

export default {
  unorderedList
};
