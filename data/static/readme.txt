复选框示例
    <div class="form-check">
        <label class="form-check-label">
          <input type="checkbox" class="form-check-input" name="needDiff" id="needDiff" value="1" checked> 显示差异
        </label>
    </div>

注：添加checked表示默认选中，必须设置name属性，id似乎没有什么用，当选中时，才会传递。表单中的数据，统一用 $('#form1').serialize() 来序列化。  

表格可用 <td  colspan="3"> 合并列， rowspan 合并行。