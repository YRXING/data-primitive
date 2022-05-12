<template>
  <div class="container">
    <el-row :gutter="20">
      <el-col :span="20">
        <div id="mountNode"></div>
      </el-col>
      <el-col :span="4" style="padding-top:10%">
        <el-row align="top">
          <el-col>
            <el-button type="primary" class="ope_btn" @click="startup">构建数源体</el-button>
          </el-col>
          <el-col>
            <el-button type="primary" class="ope_btn" @click="perception" :disabled="!isBuild">感知场景数据</el-button>
          </el-col>
          <el-col>
            <el-button type="primary" class="ope_btn" @click="cognitive" :disabled="!isPerception">认知并定义业务</el-button>
          </el-col>
          <el-col>
            <el-button type="primary" class="ope_btn" @click="linkage" :disabled="!isCognitive">联动并完成业务</el-button>
          </el-col>
        </el-row>
      </el-col>
    </el-row>
    <el-card class="process_card" v-show="processList.length !== 0">
      <div v-for="item in processList" :key="item">{{item}}</div>
    </el-card>
  </div>
</template>

<script>
import G6 from '@antv/g6';
const Util = G6.Util;
// Scale Animation
G6.registerNode(
  'circle-animate',
  {
    afterDraw(cfg, group) {
      const shape = group.get('children')[0];
      shape.animate(
        (ratio) => {
          const diff = ratio <= 0.5 ? ratio * 10 : (1 - ratio) * 10;
          return {
            r: cfg.size / 2 + diff,
          };
        },
        {
          repeat: false,
          duration: 4000,
          easing: 'easeCubic',
        },
      );
    },
  },
  'circle',
);
// Background Animation
G6.registerNode(
  'background-animate',
  {
    afterDraw(cfg, group) {
      const r = cfg.size / 2;
      const back1 = group.addShape('circle', {
        zIndex: -3,
        attrs: {
          x: 0,
          y: 0,
          r,
          fill: cfg.color,
          opacity: 0.6,
        },
        name: 'back1-shape',
      });
      const back2 = group.addShape('circle', {
        zIndex: -2,
        attrs: {
          x: 0,
          y: 0,
          r,
          fill: cfg.color,
          opacity: 0.6,
        },
        name: 'back2-shape',
      });
      const back3 = group.addShape('circle', {
        zIndex: -1,
        attrs: {
          x: 0,
          y: 0,
          r,
          fill: cfg.color,
          opacity: 0.6,
        },
        name: 'back3-shape',
      });
      group.sort(); // Sort according to the zIndex
      back1.animate(
        {
          // Magnifying and disappearing
          r: r + 10,
          opacity: 0.1,
        },
        {
          duration: 3000,
          easing: 'easeCubic',
          delay: 0,
          repeat: false, // repeat
        },
      ); // no delay
      back2.animate(
        {
          // Magnifying and disappearing
          r: r + 10,
          opacity: 0.1,
        },
        {
          duration: 3000,
          easing: 'easeCubic',
          delay: 1000,
          repeat: false, // repeat
        },
      ); // 1s delay
      back3.animate(
        {
          // Magnifying and disappearing
          r: r + 10,
          opacity: 0.1,
        },
        {
          duration: 3000,
          easing: 'easeCubic',
          delay: 2000,
          repeat: false, // repeat
        },
      ); // 3s delay
    },
  },
  'circle',
);
// Image animation
G6.registerNode(
  'inner-animate',
  {
    afterDraw(cfg, group) {
      const size = cfg.size;
      const width = size[0] - 12;
      const height = size[1] - 12;
      const image = group.addShape('image', {
        attrs: {
          x: -width / 2,
          y: -height / 2,
          width,
          height,
          img: cfg.img,
        },
        name: 'image-shape',
      });
      image.animate(
        (ratio) => {
          const toMatrix = Util.transform(
            [1, 0, 0, 0, 1, 0, 0, 0, 1],
            [['r', ratio * Math.PI * 2]],
          );
          return {
            matrix: toMatrix,
          };
        },
        {
          repeat: false,
          duration: 3000,
          easing: 'easeCubic',
        },
      );
    },
  },
  'rect',
);
const edgeTypeColorMap = {
  type1: ["#531dab", "#391085", "#391085"],
  type2: ["#d9d9d9", "#bfbfbf", "#8c8c8c"],
  type3: ["#d3adf7", "#b37feb", "#9254de"]
};

const defaultConf = {
  style: {
    lineAppendWidth: 5,
    lineDash: [0, 0],
    lineDashOffset: 0,
    opacity: 1,
    labelCfg: {
      style: {
        fillOpacity: 1
      }
    }
  },
  /**
   * 绘制边
   * @override
   * @param  {Object} cfg   边的配置项
   * @param  {G.Group} group 边的容器
   * @return {G.Shape} 图形
   */
  drawShape(cfg, group) {
    const item = group.get('item')
    const shapeStyle = this.getShapeStyle(cfg, item);
    const shape = group.addShape('path', {
      className: 'edge-path',
      attrs: shapeStyle
    });
    return shape;
  },
  drawLabel(cfg, group) {
    const labelCfg = cfg.labelCfg || {}
    const labelStyle = this.getLabelStyle(cfg, labelCfg, group)
    const text = group.addShape('text', {
      attrs: {
        ...labelStyle,
        text: cfg.label,
        fontSize: 12,
        fill: '#404040',
        cursor: 'pointer'
      },
      className: 'edge-label'
    })

    return text 
  },

  /**
   * 获取图形的配置项
   * @internal 仅在定义这一类节点使用，用户创建和更新节点
   * @param  {Object} cfg 节点的配置项
   * @return {Object} 图形的配置项
   */
  getShapeStyle(cfg, item) {
    const { startPoint, endPoint } = cfg
    const type = item.get('type')

    const defaultStyle =  this.getStateStyle('default', true, item)

    
    if(type === 'node') {
      return Object.assign({}, cfg.style, defaultStyle);
    }

    const controlPoints = this.getControlPoints(cfg);
    let points = [ startPoint ]; // 添加起始点
    // 添加控制点
    if (controlPoints) {
      points = points.concat(controlPoints);
    }
    // 添加结束点
    points.push(endPoint);
    const path = this.getPath(points);

    const style = Object.assign({}, { path }, cfg.style, defaultStyle);
    return style;
  },
  getControlPoints(cfg) {
    let controlPoints = cfg.controlPoints; // 指定controlPoints

    if (!controlPoints || !controlPoints.length) {
      const { startPoint, endPoint } = cfg;
      const innerPoint = G6.Util.getControlPoint(startPoint, endPoint, 0.5, cfg.edgeOffset || 70);
      controlPoints = [ innerPoint ];
    }
    return controlPoints;
  },
  /**
   * 获取2次贝塞尔曲线的path
   *
   * @param {array} points 起始点和两个控制点
   * @returns
   */
  getPath(points) {
    const path = [];
    path.push([ 'M', points[0].x, points[0].y ]);
    path.push([ 'Q', points[1].x, points[1].y, points[2].x, points[2].y ]);
    return path;
  },
  /**
   * 根据不同状态，获取不同状态下的样式值
   * @param {string} name 
   * @param {string} value 
   * @param {Item} item 
   */
  getStateStyle(name, value, item) {
    const model = item.getModel()
    const { style = {} } = model

    const defaultStyle = Object.assign({}, this.style)

    // 更新颜色
    return {
      ...defaultStyle,
      lineWidth: 1,
      stroke: edgeTypeColorMap[model.edgeType] && edgeTypeColorMap[model.edgeType][0],
      ...style
    }
  },
  
};

G6.registerEdge("quadratic-label-edge", defaultConf, "quadratic");
G6.registerEdge(
  'circle-running',
  {
    afterDraw(cfg, group) {
      // get the first shape in the group, it is the edge's path here=
      const shape = group.get('children')[0];
      // the start position of the edge's path
      const startPoint = shape.getPoint(0);

      // add red circle shape
      const circle = group.addShape('circle', {
        attrs: {
          x: startPoint.x,
          y: startPoint.y,
          fill: '#1890ff',
          r: 3,
        },
        name: 'circle-shape',
      });

      // animation for the red circle
      circle.animate(
        (ratio) => {
          // the operations in each frame. Ratio ranges from 0 to 1 indicating the prograss of the animation. Returns the modified configurations
          // get the position on the edge according to the ratio
          const tmpPoint = shape.getPoint(ratio);
          // returns the modified configurations here, x and y here
          return {
            x: tmpPoint.x,
            y: tmpPoint.y,
          };
        },
        {
          repeat: false, // Whether executes the animation repeatly
          duration: 2000, // the duration for executing once
        },
      );
    },
  },
  'cubic', // extend the built-in edge 'cubic'
);
G6.registerEdge(
'quadratic',
{
    afterDraw(cfg, group) {
    const shape = group.get('children')[0];
    const length = shape.getTotalLength();
    shape.animate(
        (ratio) => {
        // the operations in each frame. Ratio ranges from 0 to 1 indicating the prograss of the animation. Returns the modified configurations
        const startLen = ratio * length;
        // Calculate the lineDash
        const cfg = {
            lineDash: [startLen, length - startLen],
        };
        return cfg;
        },
        {
        repeat: false, // Whether executes the animation repeatly
        duration: 2000, // the duration for executing once
        },
    );
    },
},
'cubic', // extend the built-in edge 'cubic'
);

export default {
  data () {
    return {
      data: {
        // 数源体
        nodes: [
          {
            id: 'distributorA', // String，该节点存在则必须，节点的唯一标识
            x: 550,
            y: 120, 
            label: '经销商数原体A',
            size: 40,
            type: 'circle-animate',
            labelCfg: {
              position: 'top',
              offset: 20,
            },
            msg: '经销商数原体A'
          },
          {
            id: 'supplierA', // String，该节点存在则必须，节点的唯一标识
            x: 200,
            y: 150,
            label: '供应商数原体A',
            size: 40,
            type: 'circle-animate',
            labelCfg: {
              position: 'left',
              offset: 20,
            },
          },
          {
            id: 'supplierB', // String，该节点存在则必须，节点的唯一标识
            x: 150,
            y: 380,
            label: '供应商数原体B', 
            size: 40,
            type: 'circle-animate',
            labelCfg: {
              position: 'left',
              offset: 20,
            },
          },
          {
            id: 'bankA', // String，该节点存在则必须，节点的唯一标识
            x: 450,
            y: 520,
            label: '金融机构数原体A', 
            size: 40,
            type: 'circle-animate',
            labelCfg: {
              position: 'bottom',
              offset: 20,
            },
          },
          {
            id: 'bankB', // String，该节点存在则必须，节点的唯一标识
            x: 700,
            y: 300,
            label: '金融机构数原体B', 
            size: 40,
            type: 'circle-animate',
            labelCfg: {
              position: 'right',
              offset: 20,
            },
          },
        ],
        // 边集
        edges: [
        ],
      },
      graph: null,
      isBuild: false,
      isPerception: false,
      isCognitive: false,
      isNormal: false,
      processList: []
    }
  },
  mounted() {
    this.init()
  },
  methods: {
    init() {
      this.graph = new G6.Graph({
        container: 'mountNode', // String | HTMLElement，必须，在 Step 1 中创建的容器 id 或容器本身
        width: 1100, // Number，必须，图的宽度
        height: 800, // Number，必须，图的高度
        defaultNode: {
          style: {
            fill: '#DEE9FF',
            stroke: '#5B8FF9',
          },
        },
        defaultEdge: {
          style: {
            lineWidth: 1,
            stroke: '#b5b5b5',
          },
        },
      });
    },
    //构建数源体
    startup() {
      this.graph.data(this.data); // 读取 Step 2 中的数据源到图上
      this.graph.render(); // 渲染图
      let that = this
      that.graph.on('node:mouseenter', function (evt) {
        const node = evt.item;
        const model = node.getModel();
        model.oriLabel = model.label;
        that.graph.updateItem(node, {
          label: model.msg,
          labelCfg: {
            style: {
              fill: '#003a8c',
            },
          },
        });
      });
      that.graph.on('node:mouseleave', function (evt) {
        const node = evt.item;
        const model = node.getModel();
        that.graph.updateItem(node, {
          label: model.oriLabel,
          labelCfg: {
            style: {
              fill: '#555',
            },
          },
        });
      });
      that.graph.on('edge:mouseenter', function (evt) {
        const edge = evt.item;
        const model = edge.getModel();
        model.oriLabel = model.label;
        that.graph.updateItem(edge, {
          label: model.msg,
          labelCfg: {
            style: {
              fill: '#003a8c',
            },
          },
        });
      });
      that.graph.on('edge:mouseleave', function (evt) {
        const edge = evt.item;
        const model = edge.getModel();
        that.graph.updateItem(edge, {
          label: model.oriLabel,
          labelCfg: {
            style: {
              fill: '#555',
            },
          },
        });
      });
      setTimeout(()=>{
        that.isBuild = true
      },2000)
   },
    //感知数据
    perception() {
      this.removeItem('distributorA')
      this.graph.addItem('node',{
        id: 'distributorA', // String，该节点存在则必须，节点的唯一标识
        x: 550,
        y: 120, 
        color: '#40a9ff',
        label: '经销商数原体A',
        size: 40,
        type: 'background-animate',
        labelCfg: {
          position: 'top',
          offset: 20,
        },
        msg: 'perceived new order:&{OrderType:normal OrderPrice:10 OrderCount:10 DistributorName:distributorA'
      })
      this.addLog('distributor-log','perceived new order',600,40)
      if(!this.isNormal) {
        this.addLog('distributor-log','perceived new order',600,40)
      }else {
        this.removeItem('distributor-log')
        this.addLog('distributor-log','perceived new order',600,40)
      }
      setTimeout(()=>{
        this.isPerception = true
      },2000)
    },
    cognitive() {
      if(!this.isNormal) {
        this.removeItem('distributor-log')
        this.addLog('distributor-log','perceived new order:normal 定义业务：类型：normal,供应商进行正常生产',600,40)
        this.processList.push('类型：normal,流程：供应商进行正常生产')
        this.isCognitive = true
      }else {
        this.removeItem('distributor-log')
        this.addLog('distributor-log','perceived new order:normal 定义业务：类型：account-receivable-order,贷款审批流程',600,40)
        this.processList.push('类型：account-receivable-order,贷款审批流程')
      }
    },
    linkage() {
      if(!this.isNormal) {
        this.isNormal = true
        this.removeItem('distributor-log')
        this.addLog('distributor-log','finding supplier....',600,40)
        setTimeout(()=>{
          this.removeItem('supplierA')
          this.graph.addItem('node',{
            id: 'supplierA', // String，该节点存在则必须，节点的唯一标识
            x: 200,
            y: 150, 
            color: '#40a9ff',
            label: '供应商数原体A',
            size: 40,
            type: 'background-animate',
            labelCfg: {
              position: 'left',
              offset: 20,
            },
            msg: 'get an order {normal 10 10 distributorA}'
          })
          setTimeout(()=>{
            this.sendOrder()
          },2000)
        },3000)
      }else {
        this.sendForm()
      }
    },
    sendForm(){
      this.graph.addItem('edge',{
        id: 'distributor-to-supplier',
        source: 'distributorA',
        target: 'supplierA',
        type: 'quadratic-label-edge',
        labelCfg: {
            position: 'top',
            refY: -10,
        },
        style: {
          lineWidth: 2,
          stroke: '#e7e7e7',
          endArrow: true
        }
      })
      this.removeItem('distributor-log')
      this.addLog('distributor-log','supplier find: supplierA, establish connection successfully',600,40)
      setTimeout(()=>{
        this.removeItem('distributor-log')
        this.addLog('distributor-log','sending data to supplierA',600,40)
        this.updateLine('distributor-to-supplier', 'order', '"order_type":"normal","order_price":10,"order_count":10,"distributor_name":"distributorA"')
        setTimeout(()=>{
          this.addLog('supplier-log','get an order,start producing....',20,70)
          this.removeItem('supplier-log')
          this.addLog('supplier-log','insufficient funds,looking for a bank to make a loan..',20,70)
          this.updateLineNoType('distributor-to-supplier')
          setTimeout(()=>{
            this.removeItem('supplier-log')
            this.addLog('supplier-log','finding bank...',20,70)
            setTimeout(()=>{
              this.removeItem('bankA')
              this.graph.addItem('node',{
                id: 'bankA', // String，该节点存在则必须，节点的唯一标识
                x: 450,
                y: 520, 
                color: '#40a9ff',
                label: '金融机构数源体A',
                size: 40,
                type: 'background-animate',
                labelCfg: {
                  position: 'bottom',
                  offset: 20,
                },
                msg: 'get a form {account-receivable-order supplierA distributorA  10000} start processing.... '
              })
              setTimeout(()=>{
                this.graph.addItem('edge',{
                  id: 'supplier-to-bank',
                  source: 'supplierA',
                  target: 'bankA',
                  type: 'quadratic-label-edge',
                  labelCfg: {
                      position: 'top',
                      refY: -10,
                  },
                  style: {
                    lineWidth: 2,
                    stroke: '#e7e7e7',
                    endArrow: true
                  }
                })
                setTimeout(()=>{
                  this.removeItem('supplier-log')
                  this.addLog('supplier-log','bank find: bankA, establish connection successfully',20,70)
                  setTimeout(()=>{
                      this.removeItem('supplier-log')
                      this.addLog('supplier-log','start sending data',20,70)
                      this.updateLine('supplier-to-bank','form','Type:account-receivable-order SupplierName:supplierA DistributorName:distributorA LogisticsName: Num:10000')
                      this.addLog('bank-log','get distributor information from form...',520,520)
                      setTimeout(()=>{
                        this.updateLineNoType('supplier-to-bank')
                        this.graph.addItem('edge',{
                          id: 'bank-to-distributor',
                          source: 'bankA',
                          target: 'distributorA',
                          type: 'quadratic-label-edge',
                          labelCfg: {
                              position: 'top',
                              refY: -10,
                          },
                          style: {
                            lineWidth: 2,
                            stroke: '#e7e7e7',
                            endArrow: true
                          }
                        })
                        //经销商感知动画
                        this.addActiveDistributor('')
                        this.removeItem('supplier-log')
                        this.addLog('supplier-log','distributor find: distributorA, establish connection successfully',20,70)
                        setTimeout(()=>{
                          this.removeItem('bank-log')
                          this.addLog('bank-log','start sending data',520,520)
                          this.updateLine('bank-to-distributor','paymentPromise','"distributor_name":"distributorA","supplier_name":"supplierA","signatured":false')
                          //经销商感知动画
                          this.addActiveDistributor('get a payment promise:"distributor_name":"distributorA","supplier_name":"supplierA","signatured":false')
                          setTimeout(()=>{
                            this.removeItem('distributor-log')
                            this.addLog('distributor-log','get a payment promise',600,40)
                            this.updateLineNoType('bank-to-distributor')
                            //经销商感知动画
                            this.addActiveDistributor('')
                            setTimeout(()=>{
                              this.removeItem('distributor-log')
                              this.addLog('distributor-log','I promise to pay for products, signatured!',600,40)
                              this.graph.addItem('edge',{
                                id: 'distributor-to-bank',
                                source: 'distributorA',
                                target: 'bankA',
                                type: 'quadratic-label-edge',
                                label: 'paymentPromise',
                                msg: 'I promise to pay for products, signatured!',
                                labelCfg: {
                                    position: 'top',
                                    refY: -10,
                                },
                                style: {
                                  lineWidth: 2,
                                  stroke: 'red',
                                  endArrow: true
                                }
                              })
                              this.removeItem('bank-log')
                              this.addLog('bank-log','verify whether the payment promise is signatured...',520,520)
                              setTimeout(()=>{
                                this.updateLineNoType('distributor-to-bank')
                                this.removeItem('bank-log')
                                this.addLog('bank-log','get the payment promise from distributorA',520,520)
                                 setTimeout(()=>{
                                  this.removeItem('bank-log')
                                  this.addLog('bank-log','the loan is approved',520,520)
                                  this.graph.addItem('edge',{
                                    id: 'bank-to-supplier',
                                    source: 'bankA',
                                    target: 'supplierA',
                                    type: 'quadratic-label-edge',
                                    label: 'capital',
                                    msg: 'capital',
                                    labelCfg: {
                                        position: 'top',
                                        refY: -10,
                                    },
                                    style: {
                                      lineWidth: 2,
                                      stroke: 'red',
                                      endArrow: true
                                    }
                                  })
                                  setTimeout(()=>{
                                    this.updateLineNoType('bank-to-supplier')
                                    this.removeItem('supplier-log')
                                    this.addLog('supplier-log','products ready,start transportation...',20,70)
                                    this.graph.addItem('edge',{
                                      id: 'supplier-to-distributor',
                                      source: 'supplierA',
                                      target: 'distributorA',
                                      type: 'quadratic-label-edge',
                                      label: 'products',
                                      msg: '"Type":"account-receivable-order","supplier_name":"supplierA","distributor_name":"distributorA","logistics_name":"","num":10000',
                                      labelCfg: {
                                          position: 'top',
                                          refY: -10,
                                      },
                                      style: {
                                        lineWidth: 2,
                                        stroke: 'red',
                                        endArrow: true
                                      }
                                    })
                                    //经销商感知动画
                                    this.addActiveDistributor('prepare capital for products...')
                                    setTimeout(()=>{
                                      this.updateLineNoType('supplier-to-distributor')
                                      this.removeItem('distributor-log')
                                      this.addLog('distributor-log','prepare capital for products...',600,40)
                                      setTimeout(()=>{
                                        this.removeItem('distributor-log')
                                        //经销商感知动画
                                        this.addActiveDistributor('sending capital to bankA')
                                        this.addLog('distributor-log','sending capital to bankA',600,40)
                                        this.updateLine('distributor-to-supplier','capital','"bank_name":"bankA","num":100')
                                        setTimeout(()=>{
                                          this.updateLineNoType('distributor-to-supplier')
                                          this.removeItem('bank-log')
                                          this.addLog('bank-log','the payment result: type:TRANSPORT sourceAddress:"127.0.0.1:8082" transport:{data:"true"}',520,520)
                                          this.updateLine('bank-to-distributor','result','bankA 100')
                                          //经销商感知动画
                                          this.addActiveDistributor('')
                                          setTimeout(()=>{
                                            this.updateLineNoType('bank-to-distributor')
                                            this.removeItem('distributor-log')
                                            //经销商感知动画
                                            this.addActiveDistributor('I have received capital')
                                            this.addLog('distributor-log','I have received capital',600,40)
                                          },3000)
                                        },3000)
                                      },3000)
                                    },3000)
                                  },3000)
                                },3000)
                              },3000)
                            },3000)
                          },3000)
                        },3000)
                      },3000)
                    },3000)
                },3000)
              },3000)
            },4000)
          },2000)
        },2000)
      },3000)
    },
    sendOrder(){
      this.graph.addItem('edge',{
        id: 'distributor-to-supplier',
        source: 'distributorA',
        target: 'supplierA',
        type: 'quadratic-label-edge',
        labelCfg: {
            position: 'top',
            refY: -10,
        },
        style: {
          lineWidth: 2,
          stroke: '#e7e7e7',
          endArrow: true
        }
      })
      this.addLog('supplier-log','supplier find: supplierA, establish connection successfully',20,70)
      setTimeout(()=>{
        this.msg = 'sending data to supplierA'
        this.removeItem('distributor-log')
        this.addLog('distributor-log','sending data to supplierA',600,40)
        this.updateLine('distributor-to-supplier', 'order', '"order_type":"normal","order_price":10,"order_count":10,"distributor_name":"distributorA"')
        setTimeout(()=>{
          this.removeItem('supplier-log')
          this.addLog('supplier-log','get an order,start producing....',20,70)
          setTimeout(()=>{
            this.removeItem('supplier-log')
            this.addLog('supplier-log','products ready,start transportation...',20,70)
            this.graph.addItem('edge',{
              id: 'supplier-to-distributor',
              source: 'supplierA',
              target: 'distributorA',
              type: 'circle-running',
              label: 'products',
              labelCfg: {
                  position: 'top',
                  refY: -10,
              },
              style: {
                lineWidth: 2,
                stroke: 'red',
                endArrow: true
              },
              msg: 'normal 10 10 distributorA'
            })
          this.updateLineNoType('distributor-to-supplier','order')
          },3000)
        },2000)
      },3000)
    },
    addActiveDistributor(msg) {
      this.graph.addItem('node',{
        x: 550,
        y: 120, 
        color: '#40a9ff',
        label: '经销商数原体A',
        size: 40,
        type: 'background-animate',
        labelCfg: {
          position: 'top',
          offset: 20,
        },
        msg: msg
      })
    },
    addActiveSupplier(msg) {
      this.graph.addItem('node',{
        x: 200,
        y: 150,
        label: '供应商数原体A',
        size: 40,
        type: 'background-animate',
        labelCfg: {
          position: 'left',
          offset: 20,
        },
        msg: msg
      })
    },
    addActiveBank(msg) {
      this.graph.addItem('node',{
        x: 450,
        y: 520,
        label: '金融机构数原体A', 
        size: 40,
        type: 'circle-animate',
        labelCfg: {
          position: 'bottom',
          offset: 20,
        },
        msg: msg
      })
    },
    addLog(id,label,x,y) {
      this.graph.addItem('node',{
        id: id,
        x: x,
        y: y,
        size: [40, 40],
        type: 'inner-animate',
        img:
          'data:image/webp;base64,UklGRq4FAABXRUJQVlA4IKIFAABwHwCdASo8ADwAPiEMhEGhhv6rQAYAgS2NHsdCq/4D8AOoA64OEUAj/XPxVwyvRGvyO/gGxN/t3oK/1X6zesX3L/p/RP/2HCgKAB9AGeAbCB+AGwAbQBtA/8c/m/4PYHTonm+SzRH6B9sv2i/rOZC+G/ln9l/ML/GdoD7APcA/TD+09QDzAfrX+vHYM9AD+Uf0zrAPQA/bH0vv2t+CH9qf2R+A39e6X098/I7IAcLcs8Gjmc/9T7gPbX9H+wJ+rPVV9En9ZmBI5oUiYhkYHIVjRr9hzCTPcV5Rs/wjjIHkxPgtr/3ALZSuUm146HHwqQVA23hnnqH/4aJ/k4v4hU6RBZ0AAAD+//8ARPyL9yWIlAbWBAD0oKSqlYWreuRa3Oj02u+TvSQS8iwMYewUYTWLDNp9wOlFJaWnqE+za35UwUXuDAT6T0I4fwY+u+qrRVhl+S1ir4X7BQiNswug5AX+MjQcXEeUwfSIEUT+DFPCr+BUiwTbFxLni7fv61vRbmXoauLz4tiqOFTzEGP8tNXP5+H7mZVGfNjIxapT3FGUtqBdp/SD5cTOYOkn2fawkpqpCSqf2+CfiGWtIF673fEzlk/hIbWDhQ81C/ddxLn609d/5efckbdZ8HZbhhVmM82/Uat7CFmw1SH5xCxRxEEhjpf1EP1Xn5q9VZfm1+OFTab/MN67Xha8K//5oVBlMgZALE653X0fas/+2xMqiyCu5Wa7PHsCwbBwqROfNmzi4LPOTjkFPHVKDD1Nfj4/sul6cANdF68rf2jszlyZsUUoLTP7H3swSroc3ssNXSRVAcYd7+iBZpfoAYWvKgnr+Hv62fHZX5ZbjYbYzVjq6fsXkubto858NuUx4+ILb5y7dP6W3/IYVeUSF0yZseKIZhOMs9BBf5uB2Y3Ott//+1OG7hYINzcqigrzWAOJbSmVw3G0ULywkobx+rvfk8VmZFzQGgP/+4T74mp/vsZyM1NLguiTO2gNO05tcpXwveq5mrcweXrJ/bRZDmU2KBrnXhkXq+735c+UHTFq4h4jMOPm5shKioB6XaqhUf3DJlMg8937g/SKjrgD0H5sNm0/k4FfilBbcrsjc3dd6cwEYJo3CNhe3SpNJ2geNXyV4/hq/BZXZ1kiVknjmf5cUx7Tv/9Sb+fQ/DgAsRNSz1wiVodiLjP7aVrkxbWx5gJ8U/j0o1Ipm/nyDZRPrQXmbPAcy1eDDejxTDBKe42ElHpC2QlFdhOedsp4i9QVjt8EqWGy1YzPaGqZhCVg/LWt8/+4BmiCzNGtpR21MGJf4kI/n/1sbe36e1QBCBAx4EVfTM82ZM2lh0P189e7eY0A3NzXWVrUek8SEn+DYYCQeEaC5hDHGreFHT1baY6KyrFx3G9oMm3fLrCqmNjFRnZa3LB/5m8FgCpq9B/1OCLRE5GzVTZnVzj/4V38PgCIpX16Kznijf01+MkDeS9oCF2hEXQ9tr+mPLrjGy4Cg5fyLgyCj1fUq33nMf79Svli2h83m3gqkoxJcXvBetFQP8V/gRjBNGmFXK5TfwLhbolWEjDqUGK3n+hxzQLif9zreYO88EIRTUNbzE1/Sn7rBEtjB0uawNje5OubWsB62SOlMZoZpxrDbMb4UvQrODPhSafmhcYe9zm/dHxssMfUthhDKjyMhoRhngPjbzfGXmIV2Omgrn/zbefK/PawUGSH6x4Qk4HCN4/X8S+XCf51JJtOQeHST/yfwg69uMkE07SONnhGUrL6j5oQn6JI+zkaH/H/P/Ti/pfOTfAWxQNiMvWX08mqbuUweFSQ/G5YUP/uCvZAXutf1+Nhl2jj/n4/fPOihPjwvfFnnjOaQvs9PSpF33d+396LASZ3IID/4UP4pf9eOMXw82ccoUUUHX6MfBWyBDvARCrdPmerUwKwW+lBIAe1dsAAAA==',
        label: label,
        labelCfg: {
          position: 'right',
        },
      })
    },
    updateLine(id, label, detail) {
      const line = this.graph.findById(id);
      this.graph.updateItem(line ,{
          type: 'circle-running',
          style: {
            stroke: 'red'
          },
          label: label,
          curveOffset: [60,60],
          msg: detail
      })
    },
    updateLineNoType(id, label) {
      const line = this.graph.findById(id);
      this.graph.updateItem(line ,{
          style: {
            stroke: '#e7e7e7'
          },
          curveOffset: [60,60],
      })
    },
    removeItem(id){
      const item = this.graph.findById(id);
      this.graph.removeItem(item)
    }
  },
  beforeDestroy() {
    // this.graph.destory()
  },
  
}
</script>

<style scoped>
.container {
  width: 100%;
  height: 100%;
}
.ope_btn {
  display: flex;
  justify-content: start;
}
.ope_btn:nth-child(-n+3) {
  margin-bottom: 20px;
}
.supplier_log, .bank_log, .distributor_log{
  width: auto;
  border: #ececec 1px solid;
  position: absolute;
  text-align: left;
  border-radius: 10px;
}
.distributor_log {
  top: 0;
  right: 30%;
}
.supplier_log {
  top: -5%;
  left: 5%;
}
.bank_log {
  bottom: 20%;
  left: 60%;
}

.msg_active {
  box-shadow: 0 0 5px rgb(71, 140, 220);
}
.process_card {
  width: auto;
  position: absolute;
  right: 10%;
  bottom: 15%;
}

</style>
